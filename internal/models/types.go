package models

import "github.com/google/uuid"

type ConversationDetail struct {
	ConversationName string `json:"conversation_name"`
	ConversationID uuid.UUID `json:"id"`
	FileURL string `json:"file_url"`
	Status ConversationStatus `json:"status"`
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

type TaskStatus string

const (
	StatusTaskProcessing TaskStatus = "PROCESSING"
	StatusTaskOK         TaskStatus = "OK"
	StatusTaskError      TaskStatus = "ERROR"
)

type ConversationStatus string

const (
	StatusCreated    ConversationStatus = "CREATED"
	StatusConverted   ConversationStatus = "CONVERTED"
	StatusDiarized    ConversationStatus = "DIARIZED"
	StatusTranscribed ConversationStatus = "TRANSCRIBED"
	StatusReported    ConversationStatus = "REPORTED"
	StatusError       ConversationStatus = "ERROR"
)

type TaskType int

const (
	NoTask TaskType = iota
	ConvertTask=1
	DiarizeTask=2
	TranscribeTask=3
)