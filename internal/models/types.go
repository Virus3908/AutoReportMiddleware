package models

import "github.com/google/uuid"

type ConversationDetail struct {
	ConversationName string             `json:"conversation_name"`
	ConversationID   uuid.UUID          `json:"id"`
	FileURL          string             `json:"file_url"`
	Status           ConversationStatus `json:"status"`
	ConvertedFileURL string             `json:"converted_file_url,omitempty"`
	Segments         []SegmentDetail    `json:"segments,omitempty"`
	SemiReport       string			 	`json:"semi_report,omitempty"`
	Report           string             `json:"report,omitempty"`
}

type ConversationName struct {
	ConversationName string `json:"conversation_name"`
}

type SegmentDetail struct {
	SegmentID       uuid.UUID `json:"segment_id"`
	StartTime       float64   `json:"start_time"`
	EndTime         float64   `json:"end_time"`
	Speaker         int32     `json:"speaker"`
	ParticipantID   uuid.UUID `json:"participant_id,omitempty"`
	ParticipantName string    `json:"participant_name,omitempty"`
	TranscriptionID uuid.UUID `json:"transcription_id,omitempty"`
	Transcription   string    `json:"transcription,omitempty"`
}

type Transcription struct {
	Transcription string `json:"transcription"`
}

type Prompt struct {
	PromptName string `json:"prompt_name"`
	Prompt     string `json:"prompt"`
}

type ConnectParticipantToConversationType struct {
	ParticipantID  string `json:"participant_id"`
	ConversationID string `json:"conversation_id"`
}

type Participant struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TaskStatus int

const (
	StatusTaskError      TaskStatus = iota - 1
	StatusTaskProcessing            = 0
	StatusTaskOK                    = 1
)

type ConversationStatus int

const (
	StatusError        ConversationStatus = iota - 1
	StatusCreated                         = 0
	StatusConverted                       = 1
	StatusDiarized                        = 2
	StatusTranscribed                     = 3
	StatusSemiReported                    = 4
	StatusReported                        = 5
)

type TaskType int

const (
	NoTask         TaskType = iota
	ConvertTask             = 1
	DiarizeTask             = 2
	TranscribeTask          = 3
	SemiReportTask          = 4
	ReportTask              = 5
)


