-- +goose Up
CREATE TABLE IF NOT EXISTS rates (
    base_currency VARCHAR(3) NOT NULL,
    quote_currency VARCHAR(3) NOT NULL,
    value_e6 BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (base_currency, quote_currency)
);

CREATE TABLE IF NOT EXISTS refresh_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_currency VARCHAR(3) NOT NULL,
    quote_currency VARCHAR(3) NOT NULL,
    value_e6 BIGINT,
    status TEXT NOT NULL CHECK (status IN ('processing', 'success','error')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    error_message TEXT
);

CREATE UNIQUE INDEX rr_one_processing_per_pair
    ON refresh_requests (base_currency, quote_currency)
    WHERE status = 'processing';

-- +goose Down
DROP TABLE IF EXISTS rates;
DROP TABLE IF EXISTS refresh_requests;