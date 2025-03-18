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

func NewAPIClient(cfg APIConfig) (*APIClient, error) {
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
	_, err := getAPIResponse[ResponseInfo](context.Background(), client.Client, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create client connection: %s", err)
	}
	return &client, nil
}

func getAPIResponse[T ResponseSegments | ResponseStatus | ResponseInfo](ctx context.Context, a *http.Client, method, url string, body io.Reader) (*T, error) {
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

	response, err := getAPIResponse[ResponseStatus](ctx, a.Client, http.MethodGet, url, nil)
	if err != nil {
		return Error, fmt.Errorf("%s", err)
	}

	return response.status, nil
}

func (a *APIClient) GetDiarizationSegmentsByTaskID(ctx context.Context, ID uuid.UUID) ([]Segment, error) {
	url := a.Config.BaseURL + "/api/segments/" + ID.String()

	response, err := getAPIResponse[ResponseSegments](ctx, a.Client, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return response.Segments, nil
}

func (a *APIClient) CreateTaskConvertFileAndGetTaskID(ctx context.Context, fileURL string) (*uuid.UUID, error) {
	data := RequestFile{
		FileUrl: fileURL,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error API: %v", err)
	}

	url := a.Config.BaseURL + "/api/convert"
	response, err := getAPIResponse[ResponseStatus](ctx, a.Client, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &response.ID, nil
}
