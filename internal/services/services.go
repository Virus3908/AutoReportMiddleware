package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/internal/config"
	"net/http"
	"time"
	"github.com/google/uuid"
)


func NewAPIClient(cfg config.APIConfig) (*APIClient) {
	return &APIClient{
		BaseURL: cfg.BaseUrl,
		Client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

func getAPIResponse[T APIResponseSegments | APIResponseStatus](ctx context.Context, a *APIClient, method, url string, body io.Reader) (*T, error){
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating API request: %s", err)
	}

	resp, err := a.Client.Do(req)
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
	url := a.BaseURL + "/api/status/" + ID

	response, err := getAPIResponse[APIResponseStatus](ctx, a, http.MethodGet, url, nil)
	if err != nil {
		return Error, fmt.Errorf("%s", err)
	}

	return response.status, nil
}

func (a *APIClient) GetDiarizationSegmentsByTaskID(ctx context.Context, ID uuid.UUID) ([]Segment, error) {
	url := a.BaseURL + "/api/segments/" + ID.String()

	response, err := getAPIResponse[APIResponseSegments](ctx, a, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return response.Segments, nil
}

func (a *APIClient) CreateTaskConvertFileAndGetTaskID(ctx context.Context, fileURL string) (*uuid.UUID, error) {
	data := APIRequestFile{
		FileUrl: fileURL,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error API: %v", err)
	}

	url := a.BaseURL + "/api/convert"
	response, err := getAPIResponse[APIResponseStatus](ctx, a, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &response.ID, nil
}