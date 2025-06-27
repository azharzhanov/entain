-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    balance    NUMERIC(20, 2)        NOT NULL DEFAULT 0.0 CHECK (balance >= 0),
    created_at TIMESTAMP             NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;
-- +goose StatementEnd
