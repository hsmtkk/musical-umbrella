package hybridanalysis

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type HybridAnalysis struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

const (
	v2APIBaseURL = "https://www.hybrid-analysis.com/api/v2"
	win10x64Env  = "120"
)

func New(apiKey string) *HybridAnalysis {
	client := http.DefaultClient
	baseURL := v2APIBaseURL
	return &HybridAnalysis{client, apiKey, baseURL}
}

func NewForTest(client *http.Client, baseURL string) *HybridAnalysis {
	apiKey := "test"
	return &HybridAnalysis{client, apiKey, baseURL}
}

type SubmitFileResponse struct{}

func (h *HybridAnalysis) SubmitFile(fileName string, content []byte) (SubmitFileResponse, error) {
	contentType, reqBody, err := h.createMultiPartRequest(fileName, content)
	if err != nil {
		return SubmitFileResponse{}, err
	}
	respBytes, err := h.submitFile(contentType, reqBody)
	if err != nil {
		return SubmitFileResponse{}, err
	}
	log.Print(string(respBytes))
	return SubmitFileResponse{}, nil
}

func (h *HybridAnalysis) createMultiPartRequest(fileName string, content []byte) (string, []byte, error) {
	var buf bytes.Buffer
	multiWriter := multipart.NewWriter(&buf)
	fieldWriter, err := multiWriter.CreateFormField("environment_id")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create form; %w", err)
	}
	fieldWriter.Write([]byte(win10x64Env))
	formWriter, err := multiWriter.CreateFormFile("file", fileName)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create form; %w", err)
	}
	if _, err := io.Copy(formWriter, bytes.NewReader(content)); err != nil {
		return "", nil, fmt.Errorf("failed to write stream; %w", err)
	}
	multiWriter.Close()
	contentType := multiWriter.FormDataContentType()
	return contentType, buf.Bytes(), nil
}

func (h *HybridAnalysis) submitFile(contentType string, reqBody []byte) ([]byte, error) {
	url := h.baseURL + "/submit/file"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed make request; %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("User-Agent", "Falcon Sandbox")
	req.Header.Add("api-key", h.apiKey)
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to post; %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non 2XX HTTP status code; %d; %s", resp.StatusCode, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response; %w", err)
	}
	return body, nil
}
