from datetime import datetime
from unittest.mock import MagicMock
from zoneinfo import ZoneInfo

from crawler.fetcher import PoliteFetcher
from crawler.models import RawResponseRecord


def test_raw_saved_for_200_empty_500_403(mock_storage):
    """Verify raw responses are saved for 200, empty body, 500 error, and 403 block."""
    # 1. 200 Success
    rec1, ins1 = mock_storage.save_response(
        endpoint="/Trade/1",
        url="https://cdn.tsetmc.com/api/Trade/1",
        http_status=200,
        body='{"data": [1, 2, 3]}',
        content_type="application/json",
        ins_code="100",
    )
    assert ins1 is True
    assert rec1.id == 1
    assert rec1.http_status == 200

    # 2. Empty Body
    rec2, ins2 = mock_storage.save_response(
        endpoint="/Trade/2",
        url="https://cdn.tsetmc.com/api/Trade/2",
        http_status=200,
        body="",
        content_type="application/json",
        ins_code="100",
    )
    assert ins2 is True
    assert rec2.id == 2
    assert rec2.body == ""

    # 3. 500 Server Error
    rec3, ins3 = mock_storage.save_response(
        endpoint="/Trade/3",
        url="https://cdn.tsetmc.com/api/Trade/3",
        http_status=500,
        body="Internal Error",
        content_type="text/plain",
        ins_code="100",
    )
    assert ins3 is True
    assert rec3.http_status == 500

    # 4. 403 Forbidden / Blocked
    rec4, ins4 = mock_storage.save_response(
        endpoint="/Trade/4",
        url="https://cdn.tsetmc.com/api/Trade/4",
        http_status=403,
        body="Forbidden",
        content_type="text/html",
        ins_code="100",
        is_blocked=True,
    )
    assert ins4 is True
    assert rec4.http_status == 403
    assert rec4.is_blocked is True
    assert len(mock_storage.records) == 4


def test_idempotent_second_run(mock_storage):
    """Verify saving the same URL with identical body is idempotent (ON CONFLICT DO NOTHING)."""
    url = "https://cdn.tsetmc.com/api/Trade/GetTradeHistory/123/20261124/false"
    body = '{"trades": []}'

    rec1, ins1 = mock_storage.save_response(
        endpoint="/Trade/GetTradeHistory/123/20261124/false",
        url=url,
        http_status=200,
        body=body,
    )
    assert ins1 is True
    assert len(mock_storage.records) == 1

    # Second run with exact same URL and body
    rec2, ins2 = mock_storage.save_response(
        endpoint="/Trade/GetTradeHistory/123/20261124/false",
        url=url,
        http_status=200,
        body=body,
    )
    # Deduplicated!
    assert ins2 is False
    assert len(mock_storage.records) == 1
    assert rec2.id == rec1.id


def test_resume_skip(mock_storage, test_config):
    """Verify fetcher skips making an HTTP request if a successful response already exists."""
    url = "https://cdn.tsetmc.com/api/Trade/GetTradeHistory/123/20261124/false"
    endpoint = "/Trade/GetTradeHistory/123/20261124/false"

    # Pre-populate existing 200 response
    mock_storage.save_response(
        endpoint=endpoint,
        url=url,
        http_status=200,
        body='{"cached": true}',
    )

    mock_session = MagicMock()
    fetcher = PoliteFetcher(
        config=test_config,
        storage=mock_storage,
        session=mock_session,
    )

    outside_quiet_time = datetime(2026, 11, 24, 15, 0, 0, tzinfo=ZoneInfo("Asia/Tehran"))

    # Without refetch: skips HTTP call!
    res1 = fetcher.fetch(endpoint=endpoint, ins_code="123", refetch=False, now=outside_quiet_time)
    assert res1.is_skipped is True
    assert res1.record.body == '{"cached": true}'
    assert mock_session.get.call_count == 0

    # With refetch=True: performs HTTP call!
    mock_resp = MagicMock()
    mock_resp.status_code = 200
    mock_resp.text = '{"fresh": true}'
    mock_resp.headers = {"Content-Type": "application/json"}
    mock_session.get.return_value = mock_resp

    res2 = fetcher.fetch(endpoint=endpoint, ins_code="123", refetch=True, now=outside_quiet_time)
    assert res2.is_skipped is False
    assert mock_session.get.call_count == 1
    assert res2.record.body == '{"fresh": true}'
