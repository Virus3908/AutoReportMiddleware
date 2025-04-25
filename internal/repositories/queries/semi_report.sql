-- name: CreateSemiReport :exec
INSERT INTO
    semi_report (
        conversation_id,
        prompt_id,
        task_id,
        part_num
    )
VALUES ($1, $2, $3, $4);

-- name: UpdateSemiReportByTaskID :exec
UPDATE semi_report
SET
    semi_report = $1
WHERE
    task_id = $2;

-- name: GetCountOfUnSemiReportedParts :one
SELECT COUNT(*)
FROM
    semi_report
WHERE conversation_id = $1;