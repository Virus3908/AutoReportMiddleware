package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

type APIConfig struct {
	BaseURL string `yaml:"baseurl"`
	Timeout int    `yaml:"timeout"`
}

type APIClient struct {
	Config APIConfig
	Client *http.Client
}

type Status int

type responseInfo struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

type responseStatus struct {
	ID     uuid.UUID
	status Status
}

const (
	Queued Status = iota
	INProgress
	Done
	Error
)

type Segment struct {
	StartTime float64
	EndTime   float64
	Speaker   int32
}

type responseAudioFileSegments struct {
	ID uuid.UUID
	Segments []Segment
}

type requestFile struct {
	FileUrl string `json:"fileURL"`
}

func NewAPIClient(ctx context.Context, cfg APIConfig) (*APIClient, error) {
	if cfg.Timeout < 1 {
		return nil, fmt.Errorf("timeout can't be less 1")
	}
	client := APIClient{
		Config: cfg,
		Client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
	url := cfg.BaseURL + "/info"
	_, err := getAPIResponse[responseInfo](ctx, client.Client, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create client connection: %s", err)
	}
	return &client, nil
}

func getAPIResponse[T responseAudioFileSegments | responseStatus | responseInfo](ctx context.Context, a *http.Client, method, url string, body io.Reader) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating API request: %s", err)
	}

	resp, err := a.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error API request: %s", resp.Status)
	}

	var response T
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error API response: %s", err)
	}
	return &response, nil
}

func (a *APIClient) GetTaskStatusByID(ctx context.Context, ID string) (Status, error) {
	url := a.Config.BaseURL + "/api/status/" + ID

	response, err := getAPIResponse[responseStatus](ctx, a.Client, http.MethodGet, url, nil)
	if err != nil {
		return Error, fmt.Errorf("%s", err)
	}

	return response.status, nil
}

func (a *APIClient) GetDiarizationSegmentsByTaskID(ctx context.Context, ID uuid.UUID) ([]Segment, error) {
	url := a.Config.BaseURL + "/api/segments/" + ID.String()

	response, err := getAPIResponse[responseAudioFileSegments](ctx, a.Client, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return response.Segments, nil
}

func (a *APIClient) CreateTaskConvertFileAndGetTaskID(ctx context.Context, fileURL string) (*uuid.UUID, error) {
	data := requestFile{
		FileUrl: fileURL,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error API: %v", err)
	}

	url := a.Config.BaseURL + "/api/convert"
	response, err := getAPIResponse[responseStatus](ctx, a.Client, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &response.ID, nil
}
