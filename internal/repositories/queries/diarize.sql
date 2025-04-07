-- name: CreateDiarize :exec
INSERT INTO diarize (task_id, convert_id) VALUES ($1, $2);

-- name: GetDiarizeIDByTaskID :one
SELECT ID FROM diarize WHERE task_id = $1;