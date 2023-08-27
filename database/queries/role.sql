-- name: GetRoleById :one
SELECT
    *
FROM
    root.roles
WHERE
    id = $1
LIMIT 1;

-- name: GetRoleByName :one
SELECT
    *
FROM
    root.roles
WHERE
    name = $1
LIMIT 1;

-- name: GetAllRoles :many
SELECT
    *
FROM
    root.roles;