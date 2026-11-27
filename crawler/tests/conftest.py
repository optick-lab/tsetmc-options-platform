import hashlib
from datetime import datetime, timezone
from typing import Dict, List, Optional, Tuple

import pytest
from crawler.config import CrawlerConfig
from crawler.models import RawResponseRecord
from crawler.storage import RawStorage


class InMemoryRawStorage(RawStorage):
    """In-memory mock for RawStorage simulating PostgreSQL raw_responses table."""

    def __init__(self, config: Optional[CrawlerConfig] = None):
        super().__init__(config=config, conn=None)
        self.records: List[RawResponseRecord] = []
        self._next_id = 1
        # Map of (url, body_sha256) -> record_id for UNIQUE constraint
        self._unique_index: Dict[Tuple[str, str], int] = {}

    def find_existing_successful(self, url: str) -> Optional[RawResponseRecord]:
        for rec in reversed(self.records):
            if rec.url == url and rec.http_status == 200 and not rec.is_blocked:
                return rec
        return None

    def save_response(
        self,
        endpoint: str,
        url: str,
        http_status: int,
        body: str,
        content_type: Optional[str] = None,
        ins_code: Optional[str] = None,
        fetched_at: Optional[datetime] = None,
        is_blocked: bool = False,
    ) -> Tuple[RawResponseRecord, bool]:
        if fetched_at is None:
            fetched_at = datetime.now(timezone.utc)

        body_sha256 = hashlib.sha256(body.encode("utf-8")).hexdigest()
        key = (url, body_sha256)

        if key in self._unique_index:
            # Emulate ON CONFLICT (url, body_sha256) DO NOTHING
            existing_id = self._unique_index[key]
            for r in self.records:
                if r.id == existing_id:
                    return r, False

        rec = RawResponseRecord(
            id=self._next_id,
            endpoint=endpoint,
            url=url,
            ins_code=ins_code,
            fetched_at=fetched_at,
            http_status=http_status,
            content_type=content_type,
            body=body,
            body_sha256=body_sha256,
            is_blocked=is_blocked,
        )
        self._next_id += 1
        self.records.append(rec)
        self._unique_index[key] = rec.id
        return rec, True


@pytest.fixture
def mock_storage():
    return InMemoryRawStorage()


@pytest.fixture
def test_config():
    return CrawlerConfig(
        base_url="https://cdn.tsetmc.com/api",
        rate_limit_seconds=3.0,
        max_retries=3,
        backoff_factor=1.0,
        user_agent="tsetmc-options-platform/0.1 (student project; optick-lab)",
        request_timeout=5.0,
        quiet_hours_start="09:00",
        quiet_hours_end="12:30",
        quiet_hours_tz="Asia/Tehran",
    )
