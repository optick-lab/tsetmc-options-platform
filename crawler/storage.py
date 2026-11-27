import hashlib
import logging
from datetime import datetime, timezone
from typing import Optional, Tuple

from crawler.config import CrawlerConfig
from crawler.db import get_connection
from crawler.models import RawResponseRecord

logger = logging.getLogger(__name__)


class RawStorage:
    """Manages raw response persistence and deduplication in PostgreSQL."""

    def __init__(self, config: Optional[CrawlerConfig] = None, conn=None):
        self.config = config or CrawlerConfig.from_env()
        self._conn = conn

    def _get_connection(self):
        if self._conn is not None:
            return self._conn
        return get_connection(self.config)

    def find_existing_successful(self, url: str) -> Optional[RawResponseRecord]:
        """
        Check if a successful (HTTP 200, non-blocked) raw response already exists for the given URL.
        """
        conn = self._get_connection()
        try:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT id, endpoint, url, ins_code, fetched_at, http_status, content_type, body, body_sha256, is_blocked
                    FROM raw_responses
                    WHERE url = %s
                      AND http_status = 200
                      AND is_blocked = false
                    ORDER BY fetched_at DESC
                    LIMIT 1;
                    """,
                    (url,),
                )
                row = cur.fetchone()
                if row:
                    return RawResponseRecord(
                        id=row[0],
                        endpoint=row[1],
                        url=row[2],
                        ins_code=row[3],
                        fetched_at=row[4],
                        http_status=row[5],
                        content_type=row[6],
                        body=row[7],
                        body_sha256=row[8],
                        is_blocked=row[9],
                    )
                return None
        finally:
            if self._conn is None:
                conn.close()

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
        """
        Saves a raw HTTP response to the raw_responses table.
        Uses ON CONFLICT (url, body_sha256) DO NOTHING for idempotency.
        Returns:
            Tuple of (RawResponseRecord, bool: True if inserted, False if skipped due to existing duplicate).
        """
        if fetched_at is None:
            fetched_at = datetime.now(timezone.utc)

        body_sha256 = hashlib.sha256(body.encode("utf-8")).hexdigest()

        record = RawResponseRecord(
            id=None,
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

        conn = self._get_connection()
        inserted = False
        try:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    INSERT INTO raw_responses (
                        endpoint, url, ins_code, fetched_at, http_status, content_type, body, body_sha256, is_blocked
                    ) VALUES (
                        %s, %s, %s, %s, %s, %s, %s, %s, %s
                    )
                    ON CONFLICT (url, body_sha256) DO NOTHING
                    RETURNING id;
                    """,
                    (
                        endpoint,
                        url,
                        ins_code,
                        fetched_at,
                        http_status,
                        content_type,
                        body,
                        body_sha256,
                        is_blocked,
                    ),
                )
                row = cur.fetchone()
                if row:
                    record.id = row[0]
                    inserted = True
                conn.commit()
        except Exception as e:
            conn.rollback()
            logger.error("Failed to insert raw response into database: %s", e)
            raise
        finally:
            if self._conn is None:
                conn.close()

        return record, inserted
