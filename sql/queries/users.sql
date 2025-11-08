-- name: CreateUser :one
Insert INTO users (id, created_at, updated_at, email)
Values (
    get_random_uuid(), NOW(), NOW(), $1
)
RETURNING *;
