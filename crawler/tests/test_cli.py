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


def test_fetch_range_skips_weekends_and_prints_summary(capsys):
    """
    Verify fetch-range iterates dates, skips Thursday/Friday, and prints summary.
    Date range 20261118 (Wed) to 20261122 (Sun) contains:
      - 20261118 Wed (trading)
      - 20261119 Thu (skipped weekend)
      - 20261120 Fri (skipped weekend)
      - 20261121 Sat (trading)
      - 20261122 Sun (trading)
    Total trading days = 3.
    """
    mock_fetcher = MagicMock()

    # Day 1: fetched with trades
    res1 = FetchResult(
        record=RawResponseRecord(
            id=1,
            endpoint="/Trade/GetTradeHistory/31459957843277944/20261118/false",
            url="https://cdn.tsetmc.com/api/Trade/GetTradeHistory/31459957843277944/20261118/false",
            ins_code="31459957843277944",
            fetched_at=MagicMock(),
            http_status=200,
            content_type="application/json",
            body='{"tradeHistory": [{"nTran": 1}]}',
            body_sha256="h1",
        ),
        is_skipped=False,
        is_saved=True,
    )
    # Day 2: skipped (already stored)
    res2 = FetchResult(
        record=RawResponseRecord(
            id=2,
            endpoint="/Trade/GetTradeHistory/31459957843277944/20261121/false",
            url="https://cdn.tsetmc.com/api/Trade/GetTradeHistory/31459957843277944/20261121/false",
            ins_code="31459957843277944",
            fetched_at=MagicMock(),
            http_status=200,
            content_type="application/json",
            body='{"tradeHistory": []}',
            body_sha256="h2",
        ),
        is_skipped=True,
        is_saved=False,
    )
    # Day 3: fetched but empty trades
    res3 = FetchResult(
        record=RawResponseRecord(
            id=3,
            endpoint="/Trade/GetTradeHistory/31459957843277944/20261122/false",
            url="https://cdn.tsetmc.com/api/Trade/GetTradeHistory/31459957843277944/20261122/false",
            ins_code="31459957843277944",
            fetched_at=MagicMock(),
            http_status=200,
            content_type="application/json",
            body='{"tradeHistory": []}',
            body_sha256="h3",
        ),
        is_skipped=False,
        is_saved=True,
    )

    mock_fetcher.fetch.side_effect = [res1, res2, res3]

    exit_code = run_cli([
        "fetch-range",
        "--ins-code", "31459957843277944",
        "--from", "20261118",
        "--to", "20261122",
    ], fetcher_instance=mock_fetcher)

    assert exit_code == 0
    assert mock_fetcher.fetch.call_count == 3

    # Check that Thursday (20261119) and Friday (20261120) were never called
    called_endpoints = [call.kwargs["endpoint"] for call in mock_fetcher.fetch.call_args_list]
    assert "/Trade/GetTradeHistory/31459957843277944/20261118/false" in called_endpoints
    assert "/Trade/GetTradeHistory/31459957843277944/20261119/false" not in called_endpoints
    assert "/Trade/GetTradeHistory/31459957843277944/20261120/false" not in called_endpoints
    assert "/Trade/GetTradeHistory/31459957843277944/20261121/false" in called_endpoints
    assert "/Trade/GetTradeHistory/31459957843277944/20261122/false" in called_endpoints

    captured = capsys.readouterr().out
    assert "status: 200 (1 trades)" in captured
    assert "status: skipped (already stored)" in captured
    assert "status: empty (0 trades)" in captured
    assert "Total trading days: 3" in captured
    assert "Fetched:            1" in captured
    assert "Skipped:            1" in captured
    assert "Empty:              1" in captured
    assert "Errors:             0" in captured


def test_fetch_range_block_stops_run_immediately(capsys):
    """Verify fetch-range halts immediately upon encountering BlockedException and prints message."""
    mock_fetcher = MagicMock()
    mock_fetcher.fetch.side_effect = BlockedException("BLOCKED — stop and notify mentor")

    exit_code = run_cli([
        "fetch-range",
        "--ins-code", "31459957843277944",
        "--from", "20261101",
        "--to", "20261110",
    ], fetcher_instance=mock_fetcher)

    assert exit_code == 1
    # Stopped on first request
    assert mock_fetcher.fetch.call_count == 1
    captured = capsys.readouterr().out
    assert "BLOCKED — stop and notify mentor" in captured
