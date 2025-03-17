package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/internal/common"
	"main/internal/config"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type APIClient struct {
	BaseURL string
	Client *http.Client
}

type APIResponseStatus struct {
	ID uuid.UUID
	Status bool
	Message string
}

type APIResponseSegments struct {
	ID uuid.UUID
	StartTime float64
	EndTime float64
	Speaker int32
}

type APIRequestFile struct {
	FileUrl string `json:"fileURL"`
}

type APIRequestSegment struct {
	FileUrl string
	StartTime float64
	EndTime float64
}

func NewAPIClient(cfg config.APIConfig) (*APIClient) {
	return &APIClient{
		BaseURL: cfg.BaseUrl,
		Client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

func getResponse[T any](ctx context.Context, a *APIClient, method, url string, body io.Reader) (*T, error){
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса по API: %s", err)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("при запросе по API ошибка: %s", resp.Status)
	}

	var response T
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("ошибка ответа в API: %s", err)
	}
	return &response, nil 
}

func (a *APIClient) GetResponseStatus(ctx context.Context, ID pgtype.UUID) (*APIResponseStatus, error) {
	idStr, err := common.PGUUIDtoStr(ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации pguuid в API: %s", err)
	}
	url := a.BaseURL + "/api/status/" + idStr

	return getResponse[APIResponseStatus](ctx, a, http.MethodGet, url, nil)
}

func (a *APIClient) GetResponseSegments(ctx context.Context, ID pgtype.UUID) (*APIResponseSegments, error) {
	idStr, err := common.PGUUIDtoStr(ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации pguuid в API: %s", err)
	}
	url := a.BaseURL + "/api/segments/" + idStr

	return getResponse[APIResponseSegments](ctx, a, http.MethodGet, url, nil)
}

func (a *APIClient) ConvertFile(ctx context.Context, fileURL string) (*APIResponseStatus, error) {
	data := APIRequestFile{
		FileUrl: fileURL,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга в API: %v", err)
	}

	url := a.BaseURL + "/api/convert"
	return getResponse[APIResponseStatus](ctx, a, http.MethodPost, url, bytes.NewBuffer(jsonData))
}