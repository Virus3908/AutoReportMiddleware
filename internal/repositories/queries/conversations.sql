-- name: GetConversations :many
SELECT * FROM Conversations;

-- name: GetConversationFileURL :one
SELECT file_url FROM conversations WHERE id = $1;

-- name: DeleteConversationByID :one
DELETE FROM Conversations WHERE id = $1 RETURNING file_url;

-- name: GetConversationByID :one
SELECT * FROM Conversations WHERE id = $1;

-- name: CreateConversation :exec
INSERT INTO
    Conversations (conversation_name, file_url)
VALUES ($1, $2);

-- name: UpdateConversationStatusByConvertID :exec
UPDATE conversations
SET
    status = $1
FROM convert
WHERE
    convert.id = $2
    AND convert.conversations_id = conversations.id;

-- name: UpdateConversationStatusByDiarizeID :exec
UPDATE conversations
SET status = $1
FROM convert, diarize
WHERE
    diarize.id = $2
    AND diarize.convert_id = convert.id
    AND convert.conversations_id = conversations.id;

-- name: UpdateConversationStatusByID :exec
UPDATE conversations SET status = $1 WHERE id = $2;