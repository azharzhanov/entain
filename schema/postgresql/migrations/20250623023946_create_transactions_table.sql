-- +goose Up
-- +goose StatementBegin
CREATE TYPE TRANSACTION_STATE AS ENUM ('win', 'lose');
CREATE TYPE SOURCE_TYPE AS ENUM ('game', 'server', 'payment');

CREATE TABLE transactions
(
    id             BIGSERIAL PRIMARY KEY NOT NULL,
    transaction_id VARCHAR UNIQUE        NOT NULL,
    user_id        BIGINT                NOT NULL,
    source_type    SOURCE_TYPE           NOT NULL,
    state          TRANSACTION_STATE     NOT NULL,
    amount         NUMERIC(20, 2)        NOT NULL CHECK (amount >= 0),
    created_at     TIMESTAMP             NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS TRANSACTION_STATE;
DROP TYPE IF EXISTS SOURCE_TYPE;
-- +goose StatementEnd
