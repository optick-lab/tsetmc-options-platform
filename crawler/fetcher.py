import logging
import random
import socket
import time
from datetime import datetime, time as dtime, timezone
from typing import Callable, Optional
from zoneinfo import ZoneInfo

import requests

from crawler.config import CrawlerConfig
from crawler.exceptions import BlockedException, QuietHoursException
from crawler.models import FetchResult, RawResponseRecord
from crawler.storage import RawStorage

logger = logging.getLogger(__name__)

# DNS helper: prioritize active healthy IP for cdn.tsetmc.com when available
_orig_getaddrinfo = socket.getaddrinfo


def _tsetmc_getaddrinfo(host, port, *args, **kwargs):
    results = _orig_getaddrinfo(host, port, *args, **kwargs)
    if host == "cdn.tsetmc.com":
        healthy = [res for res in results if res[4][0] == "46.102.143.218"]
        if healthy:
            return healthy
    return results


socket.getaddrinfo = _tsetmc_getaddrinfo

BLOCK_KEYWORDS = (
    "access denied",
    "attention required",
    "cloudflare",
    "captcha",
    "security check",
    "فیلتر",
    "مسدود",
    "ip banned",
    "ddos protection",
)


def is_quiet_hours(
    now: Optional[datetime] = None,
    start_str: str = "09:00",
    end_str: str = "12:30",
    tz_str: str = "Asia/Tehran",
) -> bool:
    """
    Returns True if the current time falls within market quiet hours (09:00 - 12:30 Asia/Tehran).
    Uses Asia/Tehran timezone explicitly regardless of system or container UTC clock.
    """
    tz = ZoneInfo(tz_str)
    if now is None:
        current_dt = datetime.now(tz)
    elif now.tzinfo is None:
        # Assume naive datetime is in UTC, convert to Asia/Tehran
        current_dt = now.replace(tzinfo=timezone.utc).astimezone(tz)
    else:
        current_dt = now.astimezone(tz)

    start_h, start_m = map(int, start_str.split(":"))
    end_h, end_m = map(int, end_str.split(":"))

    quiet_start = dtime(start_h, start_m)
    quiet_end = dtime(end_h, end_m)

    return quiet_start <= current_dt.time() <= quiet_end


def is_block_response(status_code: int, body: str, content_type: Optional[str] = None) -> bool:
    """
    Detects whether an HTTP response indicates an IP ban, firewall block, or rate block.
    HTTP 403 and 429 are considered block signals.
    """
    if status_code in (403, 429):
        return True

    body_lower = body.lower() if body else ""
    if status_code in (503,) and any(k in body_lower for k in BLOCK_KEYWORDS):
        return True

    if ("<html" in body_lower or "<!doctype html" in body_lower) and any(k in body_lower for k in BLOCK_KEYWORDS):
        return True

    return False


class RateLimiter:
    """Enforces a minimum interval between requests for a single worker."""

    def __init__(
        self,
        interval_seconds: float = 300.0,
        time_fn: Callable[[], float] = time.monotonic,
        sleep_fn: Callable[[float], None] = time.sleep,
    ):
        self.interval_seconds = interval_seconds
        self.time_fn = time_fn
        self.sleep_fn = sleep_fn
        self.last_request_time: Optional[float] = None

    def wait(self) -> float:
        """Waits if necessary to respect the rate limit interval. Returns waited duration."""
        now = self.time_fn()
        waited = 0.0
        if self.last_request_time is not None:
            elapsed = now - self.last_request_time
            if elapsed < self.interval_seconds:
                sleep_duration = self.interval_seconds - elapsed
                self.sleep_fn(sleep_duration)
                waited = sleep_duration
        self.last_request_time = self.time_fn()
        return waited


class PoliteFetcher:
    """
    Polite HTTP fetcher with strict rate limiting, exponential backoff starting at >= 300s,
    honest user-agent, raw storage first, resumability, and immediate halt on block signals.
    """

    def __init__(
        self,
        config: Optional[CrawlerConfig] = None,
        storage: Optional[RawStorage] = None,
        session: Optional[requests.Session] = None,
        rate_limiter: Optional[RateLimiter] = None,
        sleep_fn: Callable[[float], None] = time.sleep,
    ):
        self.config = config or CrawlerConfig.from_env()
        self.storage = storage or RawStorage(self.config)
        self.rate_limiter = rate_limiter or RateLimiter(
            interval_seconds=self.config.min_delay_seconds,
            sleep_fn=sleep_fn,
        )
        self.sleep_fn = sleep_fn
        self.session = session or requests.Session()

    def fetch(
        self,
        endpoint: str,
        ins_code: Optional[str] = None,
        refetch: bool = False,
        force_light: bool = False,
        now: Optional[datetime] = None,
    ) -> FetchResult:
        """
        Executes a single HTTP GET request with politeness rules.
        """
        # 1. Quiet Hours check (using Asia/Tehran time)
        if is_quiet_hours(
            now=now,
            start_str=self.config.quiet_hours_start,
            end_str=self.config.quiet_hours_end,
            tz_str=self.config.quiet_hours_tz,
        ):
            if not force_light:
                msg = (
                    "Refusing to run crawler during market quiet hours "
                    f"({self.config.quiet_hours_start}-{self.config.quiet_hours_end} "
                    f"{self.config.quiet_hours_tz}). Use --force-light for single small requests."
                )
                logger.error(msg)
                raise QuietHoursException(msg)
            logger.info("Operating during quiet hours with --force-light flag for a single small request.")

        # Build full URL
        if endpoint.startswith("http://") or endpoint.startswith("https://"):
            full_url = endpoint
        else:
            full_url = f"{self.config.base_url.rstrip('/')}/{endpoint.lstrip('/')}"

        # 2. Resumability: check if successful response already exists
        if not refetch:
            existing = self.storage.find_existing_successful(full_url)
            if existing is not None:
                logger.info(
                    "Skipping fetch for %s: already successfully fetched at %s (id=%s)",
                    full_url,
                    existing.fetched_at,
                    existing.id,
                )
                return FetchResult(record=existing, is_skipped=True, is_saved=False)

        # 3. Request headers - honest User-Agent, NO spoofing, NO fake browser headers, NO cookies
        headers = {
            "User-Agent": self.config.user_agent,
            "Accept": "application/json, text/plain, */*",
        }

        # 4. Fetch loop with retry and exponential backoff
        max_attempts = self.config.max_retries + 1
        last_exception: Optional[Exception] = None
        status_code = 0
        body_text = ""
        content_type = ""
        is_blocked = False

        for attempt in range(max_attempts):
            # Enforce rate limiter before every network call
            self.rate_limiter.wait()

            fetched_at = datetime.now(timezone.utc)
            try:
                logger.info("GET %s (attempt %d/%d)", full_url, attempt + 1, max_attempts)
                resp = self.session.get(
                    full_url,
                    headers=headers,
                    timeout=self.config.request_timeout,
                    allow_redirects=True,
                )
                status_code = resp.status_code
                body_text = resp.text
                content_type = resp.headers.get("Content-Type", "")

                # Check for block FIRST: On the first block signal, stop immediately with no retry
                if is_block_response(status_code, body_text, content_type):
                    is_blocked = True
                    logger.critical("BLOCKED — stop and notify mentor")
                    # Store block response FIRST before stopping
                    self.storage.save_response(
                        endpoint=endpoint,
                        url=full_url,
                        http_status=status_code,
                        body=body_text,
                        content_type=content_type,
                        ins_code=ins_code,
                        fetched_at=fetched_at,
                        is_blocked=True,
                    )
                    raise BlockedException("BLOCKED — stop and notify mentor")

                # If transient 5xx error, retry with exponential backoff starting at >= min_delay_seconds
                if status_code in (500, 502, 503, 504) and attempt < self.config.max_retries:
                    base_delay = self.config.min_delay_seconds
                    backoff = (base_delay * (2 ** attempt)) + random.uniform(0.1, 0.5)
                    logger.warning(
                        "Received HTTP %d for %s. Retrying in %.2fs (attempt %d/%d)...",
                        status_code,
                        full_url,
                        backoff,
                        attempt + 1,
                        self.config.max_retries,
                    )
                    self.sleep_fn(backoff)
                    continue

                # Successful or non-retryable response (e.g. 200, 404, etc.)
                break

            except (requests.exceptions.Timeout, requests.exceptions.ConnectionError) as exc:
                last_exception = exc
                if attempt < self.config.max_retries:
                    base_delay = self.config.min_delay_seconds
                    backoff = (base_delay * (2 ** attempt)) + random.uniform(0.1, 0.5)
                    logger.warning(
                        "Network error (%s) on %s. Retrying in %.2fs (attempt %d/%d)...",
                        exc.__class__.__name__,
                        full_url,
                        backoff,
                        attempt + 1,
                        self.config.max_retries,
                    )
                    self.sleep_fn(backoff)
                    continue
                else:
                    logger.error("Exhausted retries for %s due to network error: %s", full_url, exc)
                    body_text = f"Network error: {exc}"
                    status_code = 0
                    break

        # 5. Raw Storage FIRST: Save response into raw_responses table
        record, inserted = self.storage.save_response(
            endpoint=endpoint,
            url=full_url,
            http_status=status_code,
            body=body_text,
            content_type=content_type,
            ins_code=ins_code,
            fetched_at=fetched_at,
            is_blocked=is_blocked,
        )

        return FetchResult(
            record=record,
            is_skipped=False,
            is_saved=inserted,
            error=str(last_exception) if last_exception and status_code == 0 else None,
        )
