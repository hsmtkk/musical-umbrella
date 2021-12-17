package virustotal_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hsmtkk/musical-umbrella/virustotal"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bs, err := os.ReadFile("upload-file.json")
		assert.Nil(t, err)
		_, err = w.Write(bs)
		assert.Nil(t, err)
	}))
	defer ts.Close()

	v := virustotal.NewForTest(ts.Client(), ts.URL)
	_, err := v.UploadFile("eicar", []byte(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`))
	assert.Nil(t, err)
}
