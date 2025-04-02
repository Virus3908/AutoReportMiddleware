-- name: GetPrompts :many
SELECT * FROM Prompts;

-- name: GetPromptByID :one
SELECT * FROM Prompts WHERE id = $1;

-- name: CreatePrompt :exec
INSERT INTO Prompts (prompt) VALUES ($1);

-- name: UpdatePromptByID :exec
UPDATE Prompts SET prompt = $1 WHERE id = $2;

-- name: DeletePromptByID :exec
DELETE FROM Prompts WHERE id = $1;