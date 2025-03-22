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

-- name: CreateConversation :exec
INSERT INTO Conversations (conversation_name, file_url) VALUES ($1, $2);

-- name: GetPromts :many
SELECT * FROM Promts;

-- name: GetPromtByID :one
SELECT * FROM Promts WHERE id = $1;

-- name: CreatePromt :exec
INSERT INTO Promts (promt) VALUES ($1);

-- name: UpdatePromtByID :exec
UPDATE Promts SET promt = $1 WHERE id = $2;

-- name: DeletePromtByID :exec
DELETE FROM Promts WHERE id = $1;

-- name: CreateConvertTask :exec
INSERT INTO Convert (conversations_id, task_id) 
VALUES ($1, $2)
ON CONFLICT (conversations_id) DO UPDATE
SET task_id = $2;


-- name: UpdateConvertTask :exec
UPDATE Convert SET file_url = $1, audio_len = $2 WHERE task_id = $3;

-- name: GetConversationFileURL :one
SELECT file_url FROM conversations WHERE id = $1;

-- name: GetConvertFileURL :one
SELECT file_url FROM convert WHERE conversations_id = $1; 

-- name: CreateDiarizeTask :exec
INSERT INTO diarize (conversation_id, task_id)
VALUES ($1, $2)
ON CONFLICT (conversation_id) DO UPDATE
SET task_id = $2;