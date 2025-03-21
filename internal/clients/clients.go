package clients

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

type Client interface {
	CreateTaskConvertFileAndGetTaskID(ctx context.Context, fileURL string) (*uuid.UUID, error)
	CreateTaskDiarizeFileAndGetTaskID(ctx context.Context, fileURL string) (*uuid.UUID, error)
	CreateTaskTranscribeSegmentFileAndGetTaskID(ctx context.Context, fileURL string, segment Segment) (*uuid.UUID, error)
	CreateTaskReportAndGetTaskID(ctx context.Context, message, promt string, audioLen float64) (*uuid.UUID, error)
	GetTaskStatusByID(ctx context.Context, ID uuid.UUID) (Status, error)
	GetDiarizationSegments(responseBody []byte) ([]Segment, error)
	GetConvertedFileURLAudioLen(responseBody []byte) (*string, *float64, error)
}

type APIConfig struct {
	BaseURL string `yaml:"baseurl"`
	Timeout int    `yaml:"timeout"`
}

type APIClient struct {
	Config APIConfig
	CallbackURL string
	Client *http.Client
}

type Status int

type responseInfo struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

type responseStatus struct {
	ID     uuid.UUID `json:"uuid"`
	Status Status    `json:"status"`
}

const (
	Queued Status = iota
	INProgress
	Done
	Error
)

type Segment struct {
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
	Speaker   int32   `json:"speaker"`
}

type responseAudioFileSegments struct {
	Task_ID  uuid.UUID `json:"id"`
	Segments []Segment `json:"segments"`
}

type responseConvertedAudioFile struct {
	FileURL string `json:"file_url"`
	AudioLen float64 `json:"audio_len"`
}

// type responseMessage struct {
// 	Message string `json:"message"`
// }



type requestFile struct {
	FileURL     string `json:"file_url"`
	CallbackURL string `json:"callback_url"`
}

type requestFileWithSegment struct {
	FileURL     string  `json:"file_url"`
	CallbackURL string  `json:"callback_url"`
	Segment     Segment `json:"segment"`
}

type requestMessageWithAudioLen struct {
	Message  string  `json:"message"`
	Promt    string  `json:"promt"`
	AudioLen float64 `json:"audio_len"`
}

func NewAPIClient(ctx context.Context, cfg APIConfig, callbackURL string) (*APIClient, error) {
	if cfg.Timeout < 1 {
		return nil, fmt.Errorf("timeout can't be less 1")
	}
	client := APIClient{
		CallbackURL: callbackURL,
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

// responseMessage | responseConvertedAudioFile | responseAudioFileSegments
func getAPIResponse[T responseStatus | responseInfo](ctx context.Context, a *http.Client, method, url string, body io.Reader) (*T, error) {
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

func (a *APIClient) CreateTaskConvertFileAndGetTaskID(ctx context.Context, fileURL string) (*uuid.UUID, error) {
	data := requestFile{
		FileURL:     fileURL,
		CallbackURL: a.CallbackURL,
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

func (a *APIClient) CreateTaskDiarizeFileAndGetTaskID(ctx context.Context, fileURL string) (*uuid.UUID, error) {
	data := requestFile{
		FileURL:     fileURL,
		CallbackURL: a.CallbackURL,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error API: %v", err)
	}

	url := a.Config.BaseURL + "/api/diarize"
	response, err := getAPIResponse[responseStatus](ctx, a.Client, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &response.ID, nil
}

func (a *APIClient) CreateTaskTranscribeSegmentFileAndGetTaskID(ctx context.Context, fileURL string, segment Segment) (*uuid.UUID, error) {
	data := requestFileWithSegment{
		FileURL:     fileURL,
		CallbackURL: a.CallbackURL,
		Segment:     segment,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error API: %v", err)
	}

	url := a.Config.BaseURL + "/api/transcribe"
	response, err := getAPIResponse[responseStatus](ctx, a.Client, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &response.ID, nil
}

func (a *APIClient) CreateTaskReportAndGetTaskID(ctx context.Context, message, promt string, audioLen float64) (*uuid.UUID, error) {
	data := requestMessageWithAudioLen{
		Message:  message,
		Promt:    promt,
		AudioLen: audioLen,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error API: %v", err)
	}

	url := a.Config.BaseURL + "/api/report"
	response, err := getAPIResponse[responseStatus](ctx, a.Client, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return &response.ID, nil
}

func (a *APIClient) GetTaskStatusByID(ctx context.Context, ID uuid.UUID) (Status, error) {
	url := a.Config.BaseURL + "/api/status/" + ID.String()

	response, err := getAPIResponse[responseStatus](ctx, a.Client, http.MethodGet, url, nil)
	if err != nil {
		return Error, fmt.Errorf("%s", err)
	}

	return response.Status, nil
}

func (a *APIClient) GetDiarizationSegments(responseBody []byte) ([]Segment, error) {
	var response responseAudioFileSegments
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("parsing backend response error %s", err)
	}
	return response.Segments, nil
}

func (a *APIClient) GetConvertedFileURLAudioLen(responseBody []byte) (*string, *float64, error){
	var response responseConvertedAudioFile
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, nil, fmt.Errorf("parsing backend response error %s", err)
	}
	return &response.FileURL, &response.AudioLen, nil
}

// func (a *APIClient) GetMessageByTaskID(ctx context.Context, ID uuid.UUID) (*string, error) {
// 	url := a.Config.BaseURL + "/api/message/" + ID.String()
// 	response, err := getAPIResponse[responseMessage](ctx, a.Client, http.MethodGet, url, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s", err)
// 	}
// 	return &response.Message, nil
// }

// func (a *APIClient) GetAudioLenByTaskID(ctx context.Context, ID uuid.UUID) (*string, *float64, error) {
// 	url := a.Config.BaseURL + "/api/converted/" + ID.String()
// 	response, err := getAPIResponse[responseConvertedAudioFile](ctx, a.Client, http.MethodGet, url, nil)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("%s", err)
// 	}
// 	return &response.FileURL, &response.AudioLen, nil
// }
