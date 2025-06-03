-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: InsertPassword :exec
UPDATE users
SET hashed_password = $1
WHERE email = $2;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;