-- name: GetAccount :one
SELECT id,
       balance
FROM accounts
WHERE id = @id;

-- name: IsAccountExists :one
SELECT EXISTS(SELECT
              FROM accounts
              WHERE id = @id) AS exists;

-- name: GetBalanceForUpdate :one
SELECT balance::FLOAT
FROM accounts
WHERE id = @id
FOR UPDATE;

-- name: UpdateBalance :exec
UPDATE accounts
SET balance = @new_balance::FLOAT
WHERE id = @id;

