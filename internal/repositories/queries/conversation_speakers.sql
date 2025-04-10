-- name: CreateSpeakerWithConversationsID :one
INSERT INTO conversation_speakers (conversation_id, speaker) VALUES ($1, $2)
RETURNING id;