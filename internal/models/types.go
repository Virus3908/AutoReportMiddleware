package models

import "github.com/google/uuid"

type ConversationDetail struct {
	ConversationName string `json:"conversation_name"`
	ConversationID uuid.UUID `json:"id"`
	FileURL string `json:"file_url"`
	Status int32 `json:"status"`
	ConvertedFileURL string `json:"converted_file_url,omitempty"`
	Segments []SegmentDetail `json:"segments,omitempty"`
}

type SegmentDetail struct {
	SegmentID uuid.UUID `json:"segment_id"`
	StartTime float64 `json:"start_time"`
	EndTime float64 `json:"end_time"`
	Speaker int32 `json:"speaker"`
	TranscriptionID uuid.UUID `json:"transcription_id,omitempty"`
	Transcription string `json:"transcription,omitempty"`
}