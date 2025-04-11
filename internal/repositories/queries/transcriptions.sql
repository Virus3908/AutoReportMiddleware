-- name: UpdateTranscriptionTextByTaskID :exec
UPDATE transcriptions SET transcription = $1 WHERE task_id = $2;

-- name: CreateTranscriptionWithTaskAndSegmentID :exec
INSERT INTO transcriptions (task_id, segment_id) VALUES ($1, $2);

-- name: GetCountOfUntranscribedSegments :one
SELECT COUNT(*)
FROM
    conversations AS c
    JOIN convert AS conv ON c.id = conv.conversations_id
    JOIN diarize AS d ON conv.id = d.convert_id
    JOIN segments AS s ON d.id = s.diarize_id
    JOIN transcriptions AS t ON s.id = t.segment_id
WHERE
    c.id = $1
    AND t.transcription IS NULL;

-- name: UpdateTranscriptionTextByID :exec
UPDATE transcriptions SET transcription = $1 WHERE id = $2;