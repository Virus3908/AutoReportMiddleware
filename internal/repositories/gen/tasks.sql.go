// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: tasks.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (task_type) VALUES ($1)
RETURNING id
`

func (q *Queries) CreateTask(ctx context.Context, taskType int32) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createTask, taskType)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteTaskByID = `-- name: DeleteTaskByID :exec
DELETE FROM tasks WHERE id = $1
`

func (q *Queries) DeleteTaskByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteTaskByID, id)
	return err
}

const getTaskByID = `-- name: GetTaskByID :one
SELECT id, status, task_type, created_at, updated_at FROM tasks WHERE id = $1
`

func (q *Queries) GetTaskByID(ctx context.Context, id uuid.UUID) (Task, error) {
	row := q.db.QueryRow(ctx, getTaskByID, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.TaskType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTasks = `-- name: GetTasks :many
SELECT id, status, task_type, created_at, updated_at FROM tasks
`

func (q *Queries) GetTasks(ctx context.Context) ([]Task, error) {
	rows, err := q.db.Query(ctx, getTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Status,
			&i.TaskType,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTaskStatus = `-- name: UpdateTaskStatus :exec
UPDATE tasks SET status = $1 WHERE id = $2
`

type UpdateTaskStatusParams struct {
	Status int32     `json:"status"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) error {
	_, err := q.db.Exec(ctx, updateTaskStatus, arg.Status, arg.ID)
	return err
}
