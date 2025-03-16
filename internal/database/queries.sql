-- name: GetUsers :many
SELECT * FROM Users;

-- name: CreateUser :exec
INSERT INTO Users (name, email) VALUES ($1, $2);

-- name: UpdateUserByID :exec
UPDATE Users SET name = $1, email = $2 WHERE id = $3;

-- name: DeleteUserByID :exec
DELETE FROM Users WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM Users WHERE id = $1;

-- name: GetConversations :many
SELECT * FROM Conversations;

-- name: DeleteConversationByID :exec
DELETE FROM Conversations WHERE id = $1;

-- name: GetConversationByID :one
SELECT * FROM Conversations WHERE id = $1;

-- name: UpdateConversationNameByID :exec
UPDATE Conversations SET conversation_name = $1 WHERE id = $2;