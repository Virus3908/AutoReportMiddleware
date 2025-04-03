-- name: GetConversations :many
SELECT * FROM Conversations;

-- name: GetConversationFileURL :one
SELECT file_url FROM conversations WHERE id = $1;

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

-- name: UpdateConversationStatusByID :exec
UPDATE conversations SET status = $1 WHERE id = $2;
