import argparse
from datetime import datetime, timedelta
import json
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

    # 1. Single Fetch
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

    # 2. Fetch Range
    range_parser = subparsers.add_parser("fetch-range", help="Fetch a date range for one instrument politely")
    range_parser.add_argument(
        "--ins-code",
        required=True,
        help="TSETMC instrument code (e.g. '31459957843277944')",
    )
    range_parser.add_argument(
        "--from",
        dest="from_date",
        required=True,
        help="Start date in YYYYMMDD format (e.g. '20261101')",
    )
    range_parser.add_argument(
        "--to",
        dest="to_date",
        required=True,
        help="End date in YYYYMMDD format (e.g. '20261123')",
    )
    range_parser.add_argument(
        "--refetch",
        action="store_true",
        default=False,
        help="Re-fetch even if a successful response already exists in raw_responses",
    )
    range_parser.add_argument(
        "--force-light",
        action="store_true",
        default=False,
        help="Permit execution during market quiet hours (09:00-12:30 Asia/Tehran)",
    )

    return parser


def run_cli(args=None, fetcher_instance: PoliteFetcher = None) -> int:
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
        fetcher = fetcher_instance or PoliteFetcher(config=config)

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

        except BlockedException:
            print("BLOCKED — stop and notify mentor")
            logger.critical("BLOCKED — stop and notify mentor")
            return 1
        except QuietHoursException as exc:
            logger.error("Quiet hours restriction: %s", exc)
            return 1
        except Exception as exc:
            logger.error("Failed to execute fetch: %s", exc)
            return 1

    elif parsed_args.command == "fetch-range":
        ins_code = parsed_args.ins_code
        try:
            start_date = datetime.strptime(parsed_args.from_date, "%Y%m%d").date()
            end_date = datetime.strptime(parsed_args.to_date, "%Y%m%d").date()
        except ValueError as err:
            logger.error("Invalid date format. Use YYYYMMDD: %s", err)
            return 2

        if start_date > end_date:
            logger.error("--from date (%s) must be <= --to date (%s)", parsed_args.from_date, parsed_args.to_date)
            return 2

        # Filter out Thursday (weekday=3) and Friday (weekday=4) (Tehran weekend)
        cur = start_date
        trading_dates = []
        while cur <= end_date:
            if cur.weekday() not in (3, 4):
                trading_dates.append(cur)
            cur += timedelta(days=1)

        config = CrawlerConfig.from_env()
        fetcher = fetcher_instance or PoliteFetcher(config=config)

        count_fetched = 0
        count_skipped = 0
        count_empty = 0
        count_errors = 0

        for d in trading_dates:
            date_str = d.strftime("%Y%m%d")
            endpoint = f"/Trade/GetTradeHistory/{ins_code}/{date_str}/false"

            try:
                result = fetcher.fetch(
                    endpoint=endpoint,
                    ins_code=ins_code,
                    refetch=parsed_args.refetch,
                    force_light=parsed_args.force_light,
                )

                if result.is_skipped:
                    count_skipped += 1
                    print(f"[{date_str}] status: skipped (already stored)")
                else:
                    rec = result.record
                    if not rec or rec.http_status != 200 or result.error:
                        count_errors += 1
                        st = rec.http_status if rec else "error"
                        print(f"[{date_str}] status: error ({st}), waiting {config.min_delay_seconds:.1f}s")
                    else:
                        # Check if tradeHistory is empty
                        body_str = (rec.body or "").strip()
                        is_empty = False
                        trade_count = 0
                        if not body_str:
                            is_empty = True
                        else:
                            try:
                                payload = json.loads(body_str)
                                trades = payload.get("tradeHistory", [])
                                if not trades:
                                    is_empty = True
                                else:
                                    trade_count = len(trades)
                            except Exception:
                                if '"tradeHistory":[]' in body_str.replace(" ", ""):
                                    is_empty = True

                        if is_empty:
                            count_empty += 1
                            print(f"[{date_str}] status: empty (0 trades), waiting {config.min_delay_seconds:.1f}s")
                        else:
                            count_fetched += 1
                            print(f"[{date_str}] status: 200 ({trade_count} trades), waiting {config.min_delay_seconds:.1f}s")

            except BlockedException:
                print("BLOCKED — stop and notify mentor")
                logger.critical("BLOCKED — stop and notify mentor")
                return 1
            except QuietHoursException as exc:
                logger.error("Quiet hours restriction: %s", exc)
                return 1
            except Exception as exc:
                logger.error("Failed on date %s: %s", date_str, exc)
                count_errors += 1

        print("==================================================")
        print("Fetch Range Summary:")
        print(f"  Instrument:         {ins_code}")
        print(f"  Date range:         {parsed_args.from_date} to {parsed_args.to_date}")
        print(f"  Total trading days: {len(trading_dates)}")
        print(f"  Fetched:            {count_fetched}")
        print(f"  Skipped:            {count_skipped}")
        print(f"  Empty:              {count_empty}")
        print(f"  Errors:             {count_errors}")
        print("==================================================")
        return 0

    return 0


def main():
    sys.exit(run_cli())
