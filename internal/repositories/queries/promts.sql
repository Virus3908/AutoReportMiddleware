-- name: GetPrompts :many
SELECT * FROM Prompts
ORDER BY created_at DESC;

-- name: GetPromptByID :one
SELECT * FROM Prompts WHERE id = $1;

-- name: CreatePrompt :exec
INSERT INTO Prompts (prompt_name, prompt) VALUES ($1, $2);

-- name: UpdatePromptByID :exec
UPDATE Prompts SET prompt = $1, prompt_name = $2 WHERE id = $3;

-- name: DeletePromptByID :exec
DELETE FROM Prompts WHERE id = $1;

-- name: GetPromptByName :one
SELECT * FROM prompts WHERE prompt_name = $1;