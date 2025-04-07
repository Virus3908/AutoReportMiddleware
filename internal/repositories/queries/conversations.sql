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
SET
    status = $1
FROM convert, diarize
WHERE
    diarize.id = $2
    AND diarize.convert_id = convert.id
    AND convert.conversations_id = conversations.id;

-- name: UpdateConversationStatusByID :exec
UPDATE conversations SET status = $1 WHERE id = $2;

-- name: GetConversationIDByTranscriptionTaskID :one
SELECT c.id
FROM
    conversations AS c
    JOIN convert AS conv ON c.id = conv.conversations_id
    JOIN diarize AS d ON conv.id = d.convert_id
    JOIN segments AS s ON d.id = s.diarize_id
    JOIN transcriptions AS t ON s.id = t.segment_id
    JOIN tasks AS tasks ON tasks.id = t.task_id
WHERE tasks.id = $1;