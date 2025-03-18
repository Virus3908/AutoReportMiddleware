package services

import (
	"github.com/google/uuid"
	"net/http"
)

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

type Status int

type APIResponseStatus struct {
	ID     uuid.UUID
	status Status
}

const (
	Queued     Status = 0
	INProgress Status = 1
	Done       Status = 2
	Error      Status = 3
)

type Segments struct {
	StartTime float64
	EndTime   float64
	Speaker   int32
}

type APIResponseSegments struct {
	ID       uuid.UUID
	Segments []Segments
}

type APIRequestFile struct {
	FileUrl string `json:"fileURL"`
}

type APIRequestSegment struct {
	FileUrl   string
	StartTime float64
	EndTime   float64
}
