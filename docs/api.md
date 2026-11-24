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
  "total": 25,
  "errors": null,
  "error_code": 0
}
```

#### Field Descriptions
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
