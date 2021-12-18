package hybridanalysis

import (
	"encoding/json"
	"fmt"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

type SubmitFileResponse struct {
	JobID         string `json:"job_id"`
	SubmissionID  string `json:"submission_id"`
	EnvironmentID int    `json:"environment_id"`
	SHA256        string `json:"sha256"`
}

func (p *Parser) SubmitFile(jsonBytes []byte) (SubmitFileResponse, error) {
	resp := SubmitFileResponse{}
	if err := json.Unmarshal(jsonBytes, &resp); err != nil {
		return resp, fmt.Errorf("failed to unmarshal JSON; %w", err)
	}
	return resp, nil
}
