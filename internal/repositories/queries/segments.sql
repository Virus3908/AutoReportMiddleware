-- name: CreateSegment :exec
INSERT INTO
    segments (
        diarize_id,
        start_time,
        end_time,
        speaker
    )
VALUES ($1, $2, $3, $4);

-- name: GetSegmentsByConversationsID :many
SELECT 
    seg.id AS segment_id, 
    seg.start_time, 
    seg.end_time, 
    conv.file_url,  
    c.id AS conversation_id
FROM
    segments AS seg
    JOIN diarize AS d ON d.id = seg.diarize_id
    JOIN convert AS conv ON conv.id = d.convert_id
    JOIN conversations AS c ON c.id = conv.conversations_id
WHERE
    c.id = $1
ORDER BY seg.start_time;