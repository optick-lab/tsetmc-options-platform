"""Custom exceptions for crawler operations."""


class CrawlerException(Exception):
    """Base exception for all crawler errors."""
    pass


class BlockedException(CrawlerException):
    """Raised when server indicates an IP or client block (403, 429, block page)."""
    pass


class QuietHoursException(CrawlerException):
    """Raised when an operation is refused during trading quiet hours without --force-light."""
    pass
