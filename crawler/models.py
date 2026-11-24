from dataclasses import dataclass
from datetime import datetime
from typing import Optional


@dataclass
class RawResponseRecord:
    id: Optional[int]
    endpoint: str
    url: str
    ins_code: Optional[str]
    fetched_at: datetime
    http_status: int
    content_type: Optional[str]
    body: str
    body_sha256: str
    is_blocked: bool = False


@dataclass
class FetchResult:
    record: Optional[RawResponseRecord]
    is_skipped: bool = False
    is_saved: bool = False
    error: Optional[str] = None
