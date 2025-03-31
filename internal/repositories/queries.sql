-- name: GetParticipants :many
SELECT * FROM Participants;

-- name: CreateParticipant :exec
INSERT INTO Participants (name, email) VALUES ($1, $2);

-- name: UpdateParticipantByID :exec
UPDATE Participants SET name = $1, email = $2 WHERE id = $3;

-- name: DeleteParticipantByID :exec
DELETE FROM Participants WHERE id = $1;

-- name: GetParticipantByID :one
SELECT * FROM Participants WHERE id = $1;

-- name: GetConversations :many
SELECT * FROM Conversations;

-- name: DeleteConversationByID :one
DELETE FROM Conversations
WHERE id = $1
RETURNING file_url;

-- name: GetConversationByID :one
SELECT * FROM Conversations WHERE id = $1;

-- name: UpdateConversationNameByID :exec
UPDATE Conversations SET conversation_name = $1 WHERE id = $2;

-- name: CreateConversation :exec
INSERT INTO Conversations (conversation_name, file_url) VALUES ($1, $2);

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

-- name: CreateConvert :exec
INSERT INTO convert (task_id) VALUES ($1);

-- name: UpdateConversationStatusByID :exec
UPDATE conversations SET status = $1 WHERE id = $2;

-- name: UpdateConvertByTaskID :exec
UPDATE convert SET file_url = $1, status = $2 WHERE task_id = $3;