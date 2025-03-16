-- name: GetConversations :many
SELECT * FROM conversations;

-- name: CreateConversation :exec
INSERT INTO conversations (conversation_name, file_url) VALUES ($1, $2);

-- name: CreateUser :exec
INSERT INTO users (name, email) VALUES ($1, $2);

-- name: GetUsers :many
SELECT * FROM users;

-- name: UpdateConversation :exec
UPDATE conversations SET status = $2 WHERE id = $1;

-- name: GetConversationByID :one
SELECT * FROM conversations WHERE id = $1;