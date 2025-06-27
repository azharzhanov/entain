-- +goose Up
-- +goose StatementBegin
INSERT INTO accounts (id, balance, created_at)
VALUES (1, 0.00, NOW()),
       (2, 0.00, NOW()),
       (3, 0.00, NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE accounts;
-- +goose StatementEnd
