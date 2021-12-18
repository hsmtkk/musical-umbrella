package hybridanalysis_test

import (
	"os"
	"testing"

	"github.com/hsmtkk/musical-umbrella/hybridanalysis"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	jsonBytes, err := os.ReadFile("./submit-file.json")
	assert.Nil(t, err)
	parsed, err := hybridanalysis.NewParser().SubmitFile(jsonBytes)
	assert.Nil(t, err)
	assert.Equal(t, "61bda3b27587d942bd5baf5d", parsed.JobID)
	assert.Equal(t, "61bda3b27587d942bd5baf5e", parsed.SubmissionID)
	assert.Equal(t, 120, parsed.EnvironmentID)
	assert.Equal(t, "8b183f5a91f61fd47774199d0b068231f6bd574a5440f172ce824a71639116f6", parsed.SHA256)
}
