# API Documentation

This document describes the REST API endpoints available in the TSETMC options platform backend service.

---

## Base URL
- Local environment: `http://localhost:8085/api/v1`
- Docker internal: `http://backend:8080/api/v1`

---

## Endpoints

### 1. Health Checks

#### `GET /api/v1/health`
Returns the operational status of the HTTP server.

**Response `200 OK`**:
```json
{
  "is_successful": true,
  "request_uuid": "c07449b2-38e9-4e89-a359-c290ee235ef9",
  "status_code": 200,
  "message": "ok",
  "data": {
    "status": "ok"
  },
  "errors": null,
  "error_code": 0
}
```

#### `GET /api/v1/health/db`
Performs a ping (`SELECT 1`) against PostgreSQL via GORM to verify live database connectivity.

**Response `200 OK`**:
```json
{
  "is_successful": true,
  "request_uuid": "9a2b069d-d36b-4ee6-8ac4-16c3a5756199",
  "status_code": 200,
  "message": "database connection healthy",
  "data": {
    "status": "ok",
    "db": "connected"
  },
  "errors": null,
  "error_code": 0
}
```

---

### 2. Underlying Assets

#### `GET /api/v1/underlyings`
Returns a paginated list of underlying instruments (stocks, ETFs, mutual funds) traded on the exchange, sorted alphabetically by symbol (`symbol ASC`).

#### Query Parameters
| Parameter | Type | Required | Default | Constraints | Description |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `page` | integer | No | `1` | `>= 1` | Page number to retrieve (1-based index). |
| `page_size` | integer | No | `20` | `1 <= page_size <= 100` | Number of items per page. |

#### Success Response (`200 OK`)
```json
{
  "is_successful": true,
  "request_uuid": "b6a8276c-38fb-41c9-a9a7-e5fa4e6720eb",
  "status_code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 20,
        "symbol": "اخابر",
        "ins_code": null,
        "asset_type": "stock"
      },
      {
        "id": 1,
        "symbol": "اطلس",
        "ins_code": null,
        "asset_type": "etf"
      }
    ],
    "page": 1,
    "page_size": 2,
    "total": 25
  },
  "errors": null,
  "error_code": 0
}
```

#### Field Descriptions
- `data`: Paginated payload object:
  - `items`: Array of underlying asset records.
    - `id`: Unique surrogate identifier (bigserial).
    - `symbol`: Official Persian symbol (e.g. `خودرو`, `همتراز`).
    - `ins_code`: TSETMC instrument code string, or `null` if not yet matched.
    - `asset_type`: Asset classification (`stock`, `etf`, `fund`).
  - `page`: The current active page number.
  - `page_size`: The number of items requested per page.
  - `total`: Total count of underlying assets across all pages.

#### Error Response (`400 Bad Request`)
Returned when invalid pagination parameters are provided (e.g. `page < 1`, `page_size < 1`, `page_size > 100`, or non-integer values).

```json
{
  "is_successful": false,
  "request_uuid": "d35f6610-1849-4702-869d-21019688ae21",
  "status_code": 400,
  "message": "page_size must be an integer greater than or equal to 1",
  "data": {},
  "errors": [
    "page_size must be an integer greater than or equal to 1"
  ],
  "error_code": 40000
}
```

---

### 3. Raw Responses & Trade Data Extraction

#### `GET /api/v1/raw-responses`
Returns the list of all archived raw responses stored in the `raw_responses` table, excluding response bodies to conserve bandwidth.

##### Success Response (`200 OK`)
```json
{
  "is_successful": true,
  "request_uuid": "4b886829-ee37-4f69-a208-66405d1115d8",
  "status_code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "endpoint": "/Trade/GetTradeHistory/31459957843277944/20261111/false",
      "ins_code": "31459957843277944",
      "http_status": 200,
      "is_blocked": false,
      "fetched_at": "2026-11-24T16:48:35.688168Z",
      "body_size": 2723
    }
  ],
  "errors": null,
  "error_code": 0
}
```

##### Field Descriptions
- `data`: Array of raw response metadata items:
  - `id`: Unique surrogate primary key.
  - `endpoint`: API endpoint path.
  - `ins_code`: TSETMC instrument code.
  - `http_status`: HTTP status code returned when fetched.
  - `is_blocked`: Flag indicating if the response was a block/captcha.
  - `fetched_at`: UTC timestamp of when response was fetched.
  - `body_size`: Byte length of the verbatim response body.

---

#### `GET /api/v1/raw-responses/:id/trades`
Parses the stored raw JSON payload of a `GetTradeHistory` response and returns individual trade items along with an aggregated summary. Excludes canceled trades (`canceled != 0`) from the summary metrics.

##### Path Parameters
- `id` (integer, required): The ID of the raw response record.

##### Success Response (`200 OK`)
```json
{
  "is_successful": true,
  "request_uuid": "8ef6f8ac-ce85-40a1-bc1d-98dfbbd825d6",
  "status_code": 200,
  "message": "success",
  "data": {
    "summary": {
      "trade_count": 14,
      "total_volume": 1251,
      "total_value": 377073,
      "min_price": 290,
      "max_price": 313,
      "first_price": 290,
      "last_price": 299
    },
    "trades": [
      {
        "n_tran": 14,
        "time": "12:29:40",
        "volume": 10,
        "price": 299,
        "canceled": false
      },
      {
        "n_tran": 13,
        "time": "12:15:21",
        "volume": 6,
        "price": 299,
        "canceled": false
      }
    ]
  },
  "errors": null,
  "error_code": 0
}
```

##### Field Descriptions
- `summary`:
  - `trade_count`: Total number of valid (non-canceled) trades.
  - `total_volume`: Sum of trade volumes for valid trades.
  - `total_value`: Sum of (volume * price) for valid trades.
  - `min_price`: Lowest transaction price among valid trades.
  - `max_price`: Highest transaction price among valid trades.
  - `first_price`: Transaction price of the earliest trade of the session.
  - `last_price`: Transaction price of the latest trade of the session.
- `trades`: Array of all trades:
  - `n_tran`: Sequential transaction number.
  - `time`: Trade timestamp formatted as `HH:MM:SS`.
  - `volume`: Trade share/contract quantity (`qTitTran`).
  - `price`: Executed transaction price (`pTran`).
  - `canceled`: Boolean flag indicating if trade was canceled.

##### Error Responses
- **`404 Not Found`**: Returned when `:id` does not exist in `raw_responses`.
- **`400 Bad Request`**: Returned when `:id` is not a positive integer, or when the endpoint of the record is not a `GetTradeHistory` response.
