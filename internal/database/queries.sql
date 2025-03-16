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

-- name: DeleteConversationByID :exec
DELETE FROM Conversations WHERE id = $1;

-- name: GetConversationByID :one
SELECT * FROM Conversations WHERE id = $1;

-- name: UpdateConversationNameByID :exec
UPDATE Conversations SET conversation_name = $1 WHERE id = $2;