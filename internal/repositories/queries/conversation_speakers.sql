-- name: CreateSpeakerWithConversationsID :one
INSERT INTO
    conversation_speakers (conversation_id, speaker)
VALUES ($1, $2)
RETURNING
    id;

-- name: AssignParticipantToSpeakerByID :exec
UPDATE conversation_speakers SET participant_id = $1 WHERE id = $2;

-- name: CreateNewSpeakerForSegment :one
INSERT INTO
    conversation_speakers (
        speaker,
        participant_id,
        conversation_id
    )
VALUES ($1, $2, $3)
RETURNING
    id;

-- name: GetSpeakerIDAndParticipantIDBySegmentID :one
SELECT cs.id as speaker_id, cs.participant_id
FROM
    conversation_speakers AS cs
    JOIN segments AS s ON cs.id = s.speaker_id
WHERE
    s.id = $1;

-- name: GetSpeakerCountByConversationID :one
SELECT COUNT(*)
FROM conversation_speakers
WHERE
    conversation_id = $1;

-- name: NullifySpeakerParticipantID :exec
UPDATE conversation_speakers SET participant_id = NULL WHERE participant_id = $1;