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

const (
	ConvertTask    = 1
	DiarizeTask    = 2
	TranscribeTask = 3
)

type Segment struct {
	Speaker int32   `json:"speaker"`
	Start   float64 `json:"start"`
	End     float64 `json:"end"`
}

type SegmentsPayload struct {
	Segments []Segment `json:"segments"`
}

type TranscriptionPayload struct {
	Text string `json:"text"`
}

const (
	StatusTaskProcessing = 1
	StatusTaskOK         = 2
	StatusTaskError      = 3
)

const (
	StatusConverted   = 1
	StatusDiarized    = 2
	StatusTranscribed = 3
	StatusReported    = 4
	StatusError       = 5
)

type Message struct {
	TaskID          uuid.UUID `json:"task_id"`
	FileURL         string    `json:"file_url"`
	StartTime       float64   `json:"start_time,omitempty"`
	EndTime         float64   `json:"end_time,omitempty"`
	CallbackPostfix string    `json:"callback_postfix"`
}