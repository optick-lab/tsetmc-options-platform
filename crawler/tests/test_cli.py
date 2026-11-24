from unittest.mock import MagicMock, patch

from crawler.cli import run_cli
from crawler.exceptions import BlockedException, QuietHoursException
from crawler.models import FetchResult, RawResponseRecord


def test_cli_fetch_template_substitution():
    """Verify CLI replaces {insCode} and {date} template parameters in --url-path."""
    with patch("crawler.cli.PoliteFetcher") as MockFetcher:
        mock_instance = MockFetcher.return_value
        mock_instance.fetch.return_value = FetchResult(
            record=RawResponseRecord(
                id=1,
                endpoint="/Trade/GetTradeHistory/46348559193224090/20261124/false",
                url="https://cdn.tsetmc.com/api/Trade/GetTradeHistory/46348559193224090/20261124/false",
                ins_code="46348559193224090",
                fetched_at=MagicMock(),
                http_status=200,
                content_type="application/json",
                body="{}",
                body_sha256="abc12345",
            ),
            is_skipped=False,
            is_saved=True,
        )

        exit_code = run_cli([
            "fetch",
            "--url-path", "/Trade/GetTradeHistory/{insCode}/{date}/false",
            "--ins-code", "46348559193224090",
            "--date", "20261124",
            "--force-light",
        ])

        assert exit_code == 0
        mock_instance.fetch.assert_called_once_with(
            endpoint="/Trade/GetTradeHistory/46348559193224090/20261124/false",
            ins_code="46348559193224090",
            refetch=False,
            force_light=True,
        )


def test_cli_blocked_exit_code():
    """Verify CLI exits with non-zero code on BlockedException."""
    with patch("crawler.cli.PoliteFetcher") as MockFetcher:
        mock_instance = MockFetcher.return_value
        mock_instance.fetch.side_effect = BlockedException("BLOCKED — stop and notify mentor")

        exit_code = run_cli([
            "fetch",
            "--url-path", "/Trade/GetTradeHistory/123/20261124/false",
            "--ins-code", "123",
            "--date", "20261124",
        ])

        assert exit_code == 1


def test_cli_quiet_hours_exit_code():
    """Verify CLI exits with non-zero code on QuietHoursException."""
    with patch("crawler.cli.PoliteFetcher") as MockFetcher:
        mock_instance = MockFetcher.return_value
        mock_instance.fetch.side_effect = QuietHoursException("Refusing to run during quiet hours")

        exit_code = run_cli([
            "fetch",
            "--url-path", "/Trade/GetTradeHistory/123/20261124/false",
            "--ins-code", "123",
            "--date", "20261124",
        ])

        assert exit_code == 1
