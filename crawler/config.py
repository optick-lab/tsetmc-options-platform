import os
from dataclasses import dataclass
from dotenv import load_dotenv

load_dotenv()


@dataclass
class CrawlerConfig:
    base_url: str = "https://cdn.tsetmc.com/api"
    min_delay_seconds: float = 300.0
    rate_limit_seconds: float = 300.0
    max_retries: int = 3
    backoff_factor: float = 1.0
    user_agent: str = "tsetmc-options-platform/0.1 (student project; optick-lab)"
    request_timeout: float = 10.0

    # Quiet hours in Asia/Tehran (09:00 - 12:30)
    quiet_hours_start: str = "09:00"
    quiet_hours_end: str = "12:30"
    quiet_hours_tz: str = "Asia/Tehran"

    # Database
    db_host: str = "localhost"
    db_port: int = 5440
    db_user: str = "postgres"
    db_password: str = "change_me"
    db_name: str = "tsetmc_options"

    def __post_init__(self):
        # Synchronize min_delay_seconds and rate_limit_seconds
        if self.min_delay_seconds != 300.0 and self.rate_limit_seconds == 300.0:
            self.rate_limit_seconds = self.min_delay_seconds
        elif self.rate_limit_seconds != 300.0 and self.min_delay_seconds == 300.0:
            self.min_delay_seconds = self.rate_limit_seconds

    @classmethod
    def from_env(cls) -> "CrawlerConfig":
        delay = float(
            os.getenv(
                "CRAWLER_MIN_DELAY_SECONDS",
                os.getenv("CRAWLER_RATE_LIMIT_SECONDS", "300.0"),
            )
        )
        return cls(
            base_url=os.getenv("CRAWLER_BASE_URL", "https://cdn.tsetmc.com/api").rstrip("/"),
            min_delay_seconds=delay,
            rate_limit_seconds=delay,
            max_retries=int(os.getenv("CRAWLER_MAX_RETRIES", "3")),
            backoff_factor=float(os.getenv("CRAWLER_BACKOFF_FACTOR", "1.0")),
            user_agent=os.getenv(
                "CRAWLER_USER_AGENT",
                "tsetmc-options-platform/0.1 (student project; optick-lab)",
            ),
            request_timeout=float(os.getenv("CRAWLER_REQUEST_TIMEOUT", "10.0")),
            quiet_hours_start=os.getenv("CRAWLER_QUIET_HOURS_START", "09:00"),
            quiet_hours_end=os.getenv("CRAWLER_QUIET_HOURS_END", "12:30"),
            quiet_hours_tz=os.getenv("CRAWLER_QUIET_HOURS_TZ", "Asia/Tehran"),
            db_host=os.getenv("DB_HOST", "localhost"),
            db_port=int(os.getenv("DB_PORT", "5440")),
            db_user=os.getenv("DB_USER", "postgres"),
            db_password=os.getenv("DB_PASSWORD", "change_me"),
            db_name=os.getenv("DB_NAME", "tsetmc_options"),
        )
