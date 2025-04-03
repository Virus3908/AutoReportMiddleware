-- name: GetTasks :many
SELECT * FROM tasks;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = $1;

-- name: CreateTask :one
INSERT INTO tasks (task_type) VALUES ($1)
RETURNING id;

-- name: UpdateTaskStatus :exec
UPDATE tasks SET status = $1 WHERE id = $2;

-- name: DeleteTaskByID :exec
DELETE FROM tasks WHERE id = $1;

