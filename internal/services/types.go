package services

import (
	"github.com/google/uuid"
	"net/http"
)

type APIClient struct {
	Config APIConfig
	Client  *http.Client
}

type Status int

type ResponseInfo struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

type ResponseStatus struct {
	ID     uuid.UUID
	status Status
}

const (
	Queued     Status = 0
	INProgress Status = 1
	Done       Status = 2
	Error      Status = 3
)

type Segment struct {
	StartTime float64
	EndTime   float64
	Speaker   int32
}

type ResponseSegments struct {
	ID       uuid.UUID
	Segments []Segment
}

type RequestFile struct {
	FileUrl string `json:"fileURL"`
}

type RequestSegment struct {
	FileUrl   string
	StartTime float64
	EndTime   float64
}
