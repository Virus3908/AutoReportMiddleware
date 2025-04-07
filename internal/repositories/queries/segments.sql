-- name: CreateSegment :exec
INSERT INTO segments (diarize_id, start_time, end_time, speaker) VALUES ($1, $2, $3, $4);