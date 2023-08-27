-- name: CreateUser :one
INSERT INTO
    root.users(username, password, role_id)
VALUES
    ($1, $2, $3)
RETURNING *;