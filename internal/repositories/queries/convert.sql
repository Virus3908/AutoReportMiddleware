-- name: GetConvert :many
SELECT * FROM convert;

-- name: GetConvertByID :one
SELECT * FROM convert WHERE id = $1;

-- name: CreateConvert :exec
INSERT INTO convert(conversations_id, task_id) VALUES ($1, $2);

-- name: UpdateConvertByTaskID :exec
UPDATE convert
SET
    file_url = $1,
    audio_len = $2
WHERE
    task_id = $3;

-- name: DeleteConvertByID :exec
DELETE FROM convert WHERE id = $1;

-- name: DeleteConvertByForgeinID :one
DELETE FROM convert WHERE conversations_id = $1 RETURNING id;

-- name: GetConvertFileURLByConversationID :one
SELECT convert.file_url, convert.ID
FROM convert
WHERE
    conversations_id = $1;