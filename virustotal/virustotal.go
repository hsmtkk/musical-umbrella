package virustotal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type VirusTotal struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

const v3APIURL = "https://www.virustotal.com/api/v3"

func New(apiKey string) *VirusTotal {
	return &VirusTotal{client: http.DefaultClient, apiKey: apiKey, baseURL: v3APIURL}
}

func NewForTest(client *http.Client, baseURL string) *VirusTotal {
	apiKey := "test"
	return &VirusTotal{client, apiKey, baseURL}
}

func (v *VirusTotal) UploadFile(fileName string, content []byte) ([]byte, error) {
	contentType, reqBody, err := v.createMultiPartRequest(fileName, content)
	if err != nil {
		return nil, err
	}
	return v.uploadFile(contentType, reqBody)
}

func (v *VirusTotal) createMultiPartRequest(fileName string, content []byte) (string, []byte, error) {
	var buf bytes.Buffer
	multiWriter := multipart.NewWriter(&buf)
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

func (v *VirusTotal) uploadFile(contentType string, reqBody []byte) ([]byte, error) {
	url := v.baseURL + "/files"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed make request; %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("x-apikey", v.apiKey)
	resp, err := v.client.Do(req)
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

func (v *VirusTotal) GetReport() {
}
