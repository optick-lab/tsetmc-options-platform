from datetime import datetime
from unittest.mock import MagicMock
from zoneinfo import ZoneInfo

import pytest
import requests

from crawler.config import CrawlerConfig
from crawler.exceptions import BlockedException, QuietHoursException
from crawler.fetcher import PoliteFetcher, RateLimiter, is_block_response, is_quiet_hours


def test_rate_limit_wait():
    """Verify RateLimiter enforces configured delay between requests using mock clock."""
    current_time = 100.0
    slept_durations = []

    def mock_time():
        return current_time

    def mock_sleep(duration):
        nonlocal current_time
        slept_durations.append(duration)
        current_time += duration

    limiter = RateLimiter(interval_seconds=3.0, time_fn=mock_time, sleep_fn=mock_sleep)

    # First request: no sleep needed
    w1 = limiter.wait()
    assert w1 == 0.0
    assert len(slept_durations) == 0

    # Advance time by only 1.0 second
    current_time = 101.0
    w2 = limiter.wait()
    # Should sleep for remaining 2.0 seconds
    assert w2 == 2.0
    assert slept_durations == [2.0]

    # Advance time by 4.0 seconds (more than interval)
    current_time = 107.0
    w3 = limiter.wait()
    assert w3 == 0.0
    assert len(slept_durations) == 1


def test_retry_backoff_on_5xx_and_timeout(mock_storage):
    """Verify exponential backoff on 5xx errors and timeouts starts at >= 300s by default."""
    slept_durations = []

    def mock_sleep(d):
        slept_durations.append(d)

    mock_session = MagicMock()
    mock_resp_500 = MagicMock()
    mock_resp_500.status_code = 500
    mock_resp_500.text = "Internal Server Error"
    mock_resp_500.headers = {"Content-Type": "text/plain"}

    mock_resp_200 = MagicMock()
    mock_resp_200.status_code = 200
    mock_resp_200.text = '{"success": true}'
    mock_resp_200.headers = {"Content-Type": "application/json"}

    # Fails twice with 500, then succeeds on 3rd attempt
    mock_session.get.side_effect = [mock_resp_500, mock_resp_500, mock_resp_200]

    # Default CrawlerConfig has min_delay_seconds = 300.0
    config = CrawlerConfig()
    rate_limiter = RateLimiter(interval_seconds=0.0)
    fetcher = PoliteFetcher(
        config=config,
        storage=mock_storage,
        session=mock_session,
        rate_limiter=rate_limiter,
        sleep_fn=mock_sleep,
    )

    outside_quiet_time = datetime(2026, 11, 24, 15, 0, 0, tzinfo=ZoneInfo("Asia/Tehran"))
    result = fetcher.fetch(
        endpoint="/Trade/GetTradeHistory/123/20261124/false",
        ins_code="123",
        now=outside_quiet_time,
    )

    assert result.record is not None
    assert result.record.http_status == 200
    assert result.record.body == '{"success": true}'
    assert mock_session.get.call_count == 3
    # Two backoff sleeps occurred, starting at >= 300s
    assert len(slept_durations) == 2
    assert slept_durations[0] >= 300.0  # 300 * 2^0 + jitter
    assert slept_durations[1] >= 600.0  # 300 * 2^1 + jitter


def test_block_detection_stops_run_and_logs_error(mock_storage, test_config, caplog):
    """Verify 403 or block page triggers BlockedException, saves response, and stops immediately."""
    mock_session = MagicMock()
    mock_resp_403 = MagicMock()
    mock_resp_403.status_code = 403
    mock_resp_403.text = "<html><title>Access Denied</title><body>Cloudflare block</body></html>"
    mock_resp_403.headers = {"Content-Type": "text/html"}

    mock_session.get.return_value = mock_resp_403

    fetcher = PoliteFetcher(
        config=test_config,
        storage=mock_storage,
        session=mock_session,
    )

    outside_quiet_time = datetime(2026, 11, 24, 15, 0, 0, tzinfo=ZoneInfo("Asia/Tehran"))

    with pytest.raises(BlockedException, match="BLOCKED — stop and notify mentor"):
        fetcher.fetch(
            endpoint="/Trade/GetTradeHistory/123/20261124/false",
            ins_code="123",
            now=outside_quiet_time,
        )

    # Verify response was stored in raw_responses before stopping
    assert len(mock_storage.records) == 1
    stored = mock_storage.records[0]
    assert stored.http_status == 403
    assert stored.is_blocked is True
    assert "Cloudflare block" in stored.body

    # Verify critical log message
    assert "BLOCKED — stop and notify mentor" in caplog.text


def test_quiet_hours_restriction(mock_storage, test_config):
    """Verify quiet hours (09:00 - 12:30 Asia/Tehran) block heavy runs unless --force-light is passed."""
    mock_session = MagicMock()
    fetcher = PoliteFetcher(
        config=test_config,
        storage=mock_storage,
        session=mock_session,
    )

    # 10:15 Tehran time is within 09:00 - 12:30
    during_quiet_time = datetime(2026, 11, 24, 10, 15, 0, tzinfo=ZoneInfo("Asia/Tehran"))

    # Without force_light: must raise QuietHoursException
    with pytest.raises(QuietHoursException):
        fetcher.fetch(
            endpoint="/Trade/GetTradeHistory/123/20261124/false",
            ins_code="123",
            force_light=False,
            now=during_quiet_time,
        )
    assert mock_session.get.call_count == 0

    # With force_light: allowed for single small request
    mock_resp = MagicMock()
    mock_resp.status_code = 200
    mock_resp.text = "{}"
    mock_resp.headers = {"Content-Type": "application/json"}
    mock_session.get.return_value = mock_resp

    res = fetcher.fetch(
        endpoint="/Trade/GetTradeHistory/123/20261124/false",
        ins_code="123",
        force_light=True,
        now=during_quiet_time,
    )
    assert res.record is not None
    assert res.record.http_status == 200
    assert mock_session.get.call_count == 1


def test_quiet_hours_uses_asia_tehran_not_utc():
    """Verify quiet-hours check uses Asia/Tehran time, not container UTC."""
    # When container UTC is 06:00, in Tehran it is 09:30 (inside quiet hours 09:00 - 12:30)
    utc_morning = datetime(2026, 11, 24, 6, 0, 0, tzinfo=ZoneInfo("UTC"))
    assert is_quiet_hours(utc_morning) is True

    # When container UTC is 09:30, in Tehran it is 13:00 (outside quiet hours)
    utc_afternoon = datetime(2026, 11, 24, 9, 30, 0, tzinfo=ZoneInfo("UTC"))
    assert is_quiet_hours(utc_afternoon) is False


def test_honest_user_agent(mock_storage, test_config):
    """Verify honest User-Agent header naming the project is sent without fake browser headers."""
    mock_session = MagicMock()
    mock_resp = MagicMock()
    mock_resp.status_code = 200
    mock_resp.text = "{}"
    mock_resp.headers = {"Content-Type": "application/json"}
    mock_session.get.return_value = mock_resp

    fetcher = PoliteFetcher(
        config=test_config,
        storage=mock_storage,
        session=mock_session,
    )

    outside_quiet_time = datetime(2026, 11, 24, 16, 0, 0, tzinfo=ZoneInfo("Asia/Tehran"))
    fetcher.fetch(
        endpoint="/Trade/GetTradeHistory/123/20261124/false",
        ins_code="123",
        now=outside_quiet_time,
    )

    args, kwargs = mock_session.get.call_args
    sent_headers = kwargs.get("headers", {})
    assert sent_headers["User-Agent"] == test_config.user_agent
    assert "tsetmc-options-platform" in sent_headers["User-Agent"]
    # No cookies or spoofed Chrome/Mozilla headers
    assert "Cookie" not in sent_headers
    assert "Sec-Ch-Ua" not in sent_headers
