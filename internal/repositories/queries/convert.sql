-- name: GetConvert :many
SELECT * FROM convert;

-- name: GetConvertByID :one
SELECT * FROM convert WHERE id = $1;

-- name: CreateConvert :exec
INSERT INTO convert (conversations_id, task_id) VALUES ($1, $2);

-- name: UpdateConvertByTaskID :exec
UPDATE convert SET file_url = $1 WHERE task_id = $2;

-- name: DeleteConvertByID :exec
DELETE FROM convert WHERE id = $1;

-- name: DeleteConvertByForgeinID :one
DELETE FROM convert 
WHERE conversations_id = $1
RETURNING id;

-- name: ASD :one
SELECT sqlc.embed(convert), sqlc.embed(conversations) FROM convert
JOIN conversations ON conversations.id = convert.conversations_id;