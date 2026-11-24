from typing import Optional
import psycopg2
from crawler.config import CrawlerConfig


def get_connection(config: Optional[CrawlerConfig] = None):
    """Establishes and returns a new PostgreSQL database connection."""
    cfg = config or CrawlerConfig.from_env()
    return psycopg2.connect(
        host=cfg.db_host,
        port=cfg.db_port,
        user=cfg.db_user,
        password=cfg.db_password,
        dbname=cfg.db_name,
    )
