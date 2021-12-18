package hybridanalysis

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
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

func (h *HybridAnalysis) SubmitFile(fileName string, content []byte) ([]byte, error) {
	contentType, reqBody, err := h.createMultiPartRequest(fileName, content)
	if err != nil {
		return nil, err
	}
	respBytes, err := h.submitFile(contentType, reqBody)
	if err != nil {
		return nil, err
	}
	log.Print(string(respBytes))
	return nil, nil
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
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": contentType,
		"User-Agent":   "Falcon Sandbox",
		"api-key":      h.apiKey,
	}
	return h.doHTTPRequest(req, headers)
}

// report endpoint needs `default` permission.
// my API key has `restricted` permission which is insufficient.
func (h *HybridAnalysis) GetReportByJobID(jobID string) ([]byte, error) {
	url := fmt.Sprintf("%s/report/%s/report/json", h.baseURL, jobID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request; %w", err)
	}
	headers := map[string]string{
		"Accept":     "application/json",
		"User-Agent": "Falcon Sandbox",
		"api-key":    h.apiKey,
	}
	return h.doHTTPRequest(req, headers)
}

func (h *HybridAnalysis) GetReportBySHA256(sha256 string, environmentID int) ([]byte, error) {
	url := fmt.Sprintf("%s/report/%s:%d/report/json", h.baseURL, sha256, environmentID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request; %w", err)
	}
	headers := map[string]string{
		"Accept":     "application/json",
		"User-Agent": "Falcon Sandbox",
		"api-key":    h.apiKey,
	}
	return h.doHTTPRequest(req, headers)
}

func (h *HybridAnalysis) doHTTPRequest(req *http.Request, headers map[string]string) ([]byte, error) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// debug
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, fmt.Errorf("failed to dump request; %w", err)
	}
	log.Print(string(reqBytes))

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request; %w", err)
	}
	defer resp.Body.Close()

	// debug
	respBytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, fmt.Errorf("failed to dump response; %w", err)
	}
	log.Print(string(respBytes))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non 2XX HTTP status code; %d; %s", resp.StatusCode, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response ; %w", err)
	}
	return body, nil
}
