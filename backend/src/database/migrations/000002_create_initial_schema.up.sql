-- Table 1: underlying_assets
CREATE TABLE IF NOT EXISTS underlying_assets (
    id BIGSERIAL PRIMARY KEY,
    symbol TEXT NOT NULL UNIQUE,
    name TEXT,
    ins_code TEXT UNIQUE,
    asset_type TEXT NOT NULL CHECK (asset_type IN ('stock', 'etf', 'fund')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Table 2: option_contracts
CREATE TABLE IF NOT EXISTS option_contracts (
    id BIGSERIAL PRIMARY KEY,
    ins_code TEXT NOT NULL UNIQUE,
    symbol TEXT NOT NULL,
    underlying_id BIGINT NOT NULL REFERENCES underlying_assets(id) ON DELETE RESTRICT,
    option_type TEXT NOT NULL CHECK (option_type IN ('call', 'put')),
    strike_price NUMERIC NOT NULL,
    begin_date DATE,
    end_date DATE NOT NULL,
    contract_size INTEGER NOT NULL DEFAULT 1000,
    isin TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_option_contracts_underlying_id ON option_contracts(underlying_id);
CREATE INDEX IF NOT EXISTS idx_option_contracts_end_date ON option_contracts(end_date);

-- Table 3: raw_responses
CREATE TABLE IF NOT EXISTS raw_responses (
    id BIGSERIAL PRIMARY KEY,
    endpoint TEXT NOT NULL,
    url TEXT NOT NULL,
    ins_code TEXT,
    fetched_at TIMESTAMPTZ NOT NULL,
    http_status INTEGER NOT NULL,
    content_type TEXT,
    body TEXT NOT NULL,
    body_sha256 TEXT NOT NULL,
    is_blocked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_raw_responses_url_body_sha256 UNIQUE (url, body_sha256)
);

CREATE INDEX IF NOT EXISTS idx_raw_responses_endpoint_ins_code_fetched_at ON raw_responses(endpoint, ins_code, fetched_at);
