-- name: CreateTransaction :exec
INSERT INTO transactions (transaction_id, user_id, source_type, state, amount, created_at)
VALUES (@transaction_id, @user_id, @source_type, @state, @amount, NOW());

-- name: IsTransactionExists :one
SELECT EXISTS (
    SELECT
    FROM transactions
    WHERE transaction_id = @transaction_id
) AS exists;