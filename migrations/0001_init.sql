-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS rates (
    base_currency VARCHAR(3) NOT NULL,
    quote_currency VARCHAR(3) NOT NULL,
    value_e6 BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT rates_base_quote_different CHECK (base_currency <> quote_currency),
    CONSTRAINT rates_currency_uppercase CHECK (
        base_currency = UPPER(base_currency) AND quote_currency = UPPER(quote_currency)
    ),

    PRIMARY KEY (base_currency, quote_currency)
);

CREATE TABLE IF NOT EXISTS refresh_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_currency VARCHAR(3) NOT NULL,
    quote_currency VARCHAR(3) NOT NULL,
    value_e6 BIGINT,
    status TEXT NOT NULL CHECK (status IN ('processing', 'success','error')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    error_message TEXT,

    CONSTRAINT rr_base_quote_different CHECK (base_currency <> quote_currency),
    CONSTRAINT rr_currency_uppercase CHECK (
        base_currency = UPPER(base_currency) AND quote_currency = UPPER(quote_currency)
    ),
    CONSTRAINT rr_success_requires_value CHECK (status <> 'success' OR value_e6 IS NOT NULL),
    CONSTRAINT rr_error_requires_message CHECK (status <> 'error' OR error_message IS NOT NULL)
);

CREATE UNIQUE INDEX rr_one_processing_per_pair
    ON refresh_requests (base_currency, quote_currency)
    WHERE status = 'processing';

-- +goose Down
DROP TABLE IF EXISTS rates;
DROP TABLE IF EXISTS refresh_requests;