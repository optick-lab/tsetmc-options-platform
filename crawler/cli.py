import argparse
import logging
import sys

from crawler.config import CrawlerConfig
from crawler.exceptions import BlockedException, QuietHoursException
from crawler.fetcher import PoliteFetcher

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger("crawler")


def create_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="crawler",
        description="TSETMC Options Platform - Polite Data Crawler",
    )
    subparsers = parser.add_subparsers(dest="command", required=True)

    fetch_parser = subparsers.add_parser("fetch", help="Fetch a single endpoint politely")
    fetch_parser.add_argument(
        "--url-path",
        required=True,
        help="URL path template (e.g. '/Trade/GetTradeHistory/{insCode}/{date}/false')",
    )
    fetch_parser.add_argument(
        "--ins-code",
        default=None,
        help="TSETMC instrument code (e.g. '46348559193224090')",
    )
    fetch_parser.add_argument(
        "--date",
        default=None,
        help="Date parameter in YYYYMMDD format",
    )
    fetch_parser.add_argument(
        "--refetch",
        action="store_true",
        default=False,
        help="Re-fetch even if a successful response already exists in raw_responses",
    )
    fetch_parser.add_argument(
        "--force-light",
        action="store_true",
        default=False,
        help="Permit single small request during market quiet hours (09:00-12:30 Asia/Tehran)",
    )

    return parser


def run_cli(args=None) -> int:
    parser = create_parser()
    parsed_args = parser.parse_args(args)

    if parsed_args.command == "fetch":
        endpoint = parsed_args.url_path
        if "{insCode}" in endpoint:
            if not parsed_args.ins_code:
                logger.error("--ins-code is required when --url-path contains {insCode}")
                return 2
            endpoint = endpoint.replace("{insCode}", parsed_args.ins_code)

        if "{date}" in endpoint:
            if not parsed_args.date:
                logger.error("--date is required when --url-path contains {date}")
                return 2
            endpoint = endpoint.replace("{date}", parsed_args.date)

        config = CrawlerConfig.from_env()
        fetcher = PoliteFetcher(config=config)

        try:
            result = fetcher.fetch(
                endpoint=endpoint,
                ins_code=parsed_args.ins_code,
                refetch=parsed_args.refetch,
                force_light=parsed_args.force_light,
            )

            if result.is_skipped:
                logger.info(
                    "Skipped: URL %s already exists in raw_responses (status=%d, sha256=%s...)",
                    result.record.url if result.record else endpoint,
                    result.record.http_status if result.record else 200,
                    result.record.body_sha256[:8] if result.record else "",
                )
                return 0

            rec = result.record
            if rec:
                logger.info(
                    "Fetched %s: HTTP %d (bytes=%d, sha256=%s, saved_id=%s)",
                    rec.url,
                    rec.http_status,
                    len(rec.body),
                    rec.body_sha256[:8],
                    rec.id,
                )
            return 0

        except BlockedException as exc:
            logger.critical("Fatal: %s", exc)
            return 1
        except QuietHoursException as exc:
            logger.error("Quiet hours restriction: %s", exc)
            return 1
        except Exception as exc:
            logger.error("Failed to execute fetch: %s", exc)
            return 1

    return 0


def main():
    sys.exit(run_cli())
