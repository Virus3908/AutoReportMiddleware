-- name: GetConvert :many
SELECT * FROM convert;

-- name: GetConvertByID :one
SELECT * FROM convert WHERE id = $1;

-- name: CreateConvert :one
INSERT INTO convert (conversations_id) VALUES ($1)
RETURNING id;

-- name: UpdateConvertByTaskID :exec
UPDATE convert SET file_url = $1, status = $2 WHERE id = $3;

-- name: DeleteConvertByID :exec
DELETE FROM convert WHERE id = $1;

-- name: DeleteConvertByForgeinID :one
DELETE FROM convert 
WHERE conversations_id = $1
RETURNING id;