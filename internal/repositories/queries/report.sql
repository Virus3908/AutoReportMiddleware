-- name: CreateReport :exec
INSERT INTO
    reports (
        conversation_id,
        prompt_id,
        task_id
    )
VALUES ($1, $2, $3);

-- name: UpdateReportByTaskID :exec
UPDATE reports
SET
    report = $1
WHERE
    task_id = $2;

-- name: GetReportByConversationID :one
SELECT *
FROM
    reports
WHERE conversation_id = $1
ORDER BY updated_at DESC;