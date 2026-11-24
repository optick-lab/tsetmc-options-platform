# Crawler Architecture and Politeness Guide

This document outlines the design, politeness guardrails, raw storage architecture, and CLI usage for the `crawler` module in `tsetmc-options-platform`.

---

## 1. Overview

The crawler is responsible for responsibly fetching financial data from the TSETMC CDN (`https://cdn.tsetmc.com/api/...`) and storing the exact raw HTTP response in PostgreSQL (`raw_responses` table) before any downstream parsing occurs.

### Key Principles
- **Polite & Respectful**: Conservative request pacing, zero concurrency, honest identification.
- **Raw-First Storage**: Store verbatim HTTP responses (success, empty, server error, or block) prior to parsing.
- **Idempotent & Resumable**: Skips already-fetched endpoints, ignores duplicate payloads via content-hash constraints.
- **Fail-Safe Block Detection**: Immediately halts on access restrictions without attempting bypasses.

---

## 2. Politeness Guardrails

### 2.1 Fixed Slow Rate Limit (Single Worker)
- **Rate**: 1 request every 3.0 seconds by default (configurable via `CRAWLER_RATE_LIMIT_SECONDS`).
- **Concurrency**: Strictly single worker (`concurrency = 1`), no parallel threads or async workers.
- **Enforcement**: In-memory rate limiter tracks timestamps between consecutive outbound requests and sleeps as necessary.

### 2.2 Honest User-Agent
- **Header**: Identified as `tsetmc-options-platform/0.1 (student project; optick-lab)` (configurable via `CRAWLER_USER_AGENT`).
- **Rules**:
  - No spoofing or browser user-agent mimicking.
  - No fake headers (e.g. spoofed `Referer`, `Sec-Ch-Ua`).
  - No proxies, VPN rotation, or cookie harvesting.

### 2.3 Exponential Backoff with Jitter
- **Triggers**: Network timeouts, HTTP 5xx errors, or HTTP 429 (Too Many Requests).
- **Formula**: `sleep_time = (backoff_factor * 2^attempt) + random_jitter(0, 1)`.
- **Max Retries**: Default 3 retries (configurable via `CRAWLER_MAX_RETRIES`).
- If retries are exhausted, the failure response or error state is logged and stored.

### 2.4 Block Detection & Immediate Halt
- **Triggers**:
  - HTTP status `403 Forbidden` or repeated `429 Too Many Requests`.
  - Body text containing signature block indicators (e.g. `"cloudflare"`, `"access denied"`, `"ip blocked"`, `"captcha"`, `"waf"`).
- **Behavior**:
  - Saves the blocked response record to `raw_responses` with `is_blocked = true`.
  - Emits an urgent log: `BLOCKED — stop and notify mentor`.
  - Immediately aborts the execution with a non-zero exit code (`exit 1`).
  - **Never attempts evasion, IP hopping, or header mutation.**

### 2.5 Market Quiet Hours (Asia/Tehran)
- **Peak Window**: `09:00` to `12:30` Asia/Tehran timezone (`UTC+03:30`), during active Tehran Stock Exchange trading sessions.
- **Enforcement**: Bulk or automated runs are rejected with `QuietHoursException` during this window.
- **Override Flag**: `--force-light` allows running a single isolated request for development/debugging inspection.

---

## 3. Raw Storage First & Idempotency

### 3.1 Raw Storage (`raw_responses` Table)
Before any data normalization or extraction takes place, the crawler persists the response in `raw_responses`:
- `endpoint`: Normalized request endpoint path (e.g. `/Trade/GetTradeHistory/...`).
- `url`: Full target URL.
- `ins_code`: TSETMC instrument code (if applicable).
- `http_status`: HTTP status code returned (200, 403, 500, etc.).
- `content_type`: `Content-Type` header from response.
- `body`: Verbatim response body (stored as text/raw bytes).
- `body_sha256`: SHA-256 hash of the response body.
- `is_blocked`: Boolean flag indicating if the response triggered block heuristics.
- `fetched_at`: UTC timestamp of the request.

### 3.2 Idempotency
- The database enforces `UNIQUE (url, body_sha256)`.
- Insertion uses `ON CONFLICT (url, body_sha256) DO NOTHING`.
- If a fetch produces a payload already saved for that exact URL, the database silently preserves the existing row and avoids duplicates.

### 3.3 Resumability
- Before issuing a network request, the fetcher queries `raw_responses` for an existing successful (`http_status = 200` AND `is_blocked = false`) entry for the requested URL.
- If present, the network call is skipped entirely, returning the cached record.
- Use `--refetch` to force a new network request when needed.

---

## 4. Environment Configuration

All crawler options are loaded from environment variables (see `.env.example`):

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` / `5440` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `change_me` | Database password |
| `DB_NAME` | `tsetmc_options` | Database name |
| `CRAWLER_BASE_URL` | `https://cdn.tsetmc.com/api` | Base API URL |
| `CRAWLER_RATE_LIMIT_SECONDS` | `3.0` | Minimum wait between requests |
| `CRAWLER_MAX_RETRIES` | `3` | Maximum retry attempts |
| `CRAWLER_BACKOFF_FACTOR` | `1.0` | Base exponential backoff multiplier |
| `CRAWLER_USER_AGENT` | `tsetmc-options-platform/0.1 (student project; optick-lab)` | Honest User-Agent |
| `CRAWLER_REQUEST_TIMEOUT` | `10.0` | Request timeout in seconds |
| `CRAWLER_QUIET_HOURS_START` | `09:00` | Start of peak trading hours |
| `CRAWLER_QUIET_HOURS_END` | `12:30` | End of peak trading hours |
| `CRAWLER_QUIET_HOURS_TZ` | `Asia/Tehran` | Timezone for quiet hours |

---

## 5. CLI Usage

The CLI handles single-endpoint fetches with parameter interpolation.

### 5.1 Basic Fetch Command
```bash
python -m crawler fetch \
  --url-path "/Trade/GetTradeHistory/{insCode}/{date}/false" \
  --ins-code "46348559193224090" \
  --date "20241124"
```

### 5.2 Command Flags
- `--url-path`: Relative endpoint path. `{insCode}` and `{date}` are replaced if provided.
- `--ins-code`: Instrument code for interpolation and record tracking.
- `--date`: Date parameter (`YYYYMMDD`).
- `--refetch`: Bypass resume check and re-fetch from TSETMC even if already cached.
- `--force-light`: Permit a single light request during quiet hours (`09:00 - 12:30 Asia/Tehran`).

### 5.3 Docker Compose Execution
The crawler is integrated as a one-shot container in `docker-compose.yml`:

```bash
# View help
docker compose run --rm crawler python -m crawler --help

# Run a single fetch (with --force-light if during market hours)
docker compose run --rm crawler python -m crawler fetch \
  --url-path "/Trade/GetTradeHistory/{insCode}/{date}/false" \
  --ins-code "46348559193224090" \
  --date "20241124" \
  --force-light
```

---

## 6. Testing

All unit tests strictly mock HTTP requests and database interactions; no outbound network traffic is generated.

```bash
# Run crawler test suite
PYTHONPATH=. pytest -v crawler/tests
```
