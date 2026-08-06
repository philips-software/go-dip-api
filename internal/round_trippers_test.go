package internal_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/philips-software/go-dip-api/internal"
	"github.com/stretchr/testify/assert"
)

type dummyTransport struct {
	responseHeader http.Header
	responseBody   string
}

func (d *dummyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	if d.responseHeader != nil {
		for k, vv := range d.responseHeader {
			for _, v := range vv {
				rec.Header().Add(k, v)
			}
		}
	}
	rec.WriteHeader(http.StatusOK)
	body := d.responseBody
	if body == "" {
		body = `{"status":"ok"}`
	}
	_, _ = rec.WriteString(body)
	return rec.Result(), nil
}

func TestLoggingRoundTripper_Redaction(t *testing.T) {
	t.Run("Hsdp-Api-Signature request header redaction", func(t *testing.T) {
		var buf bytes.Buffer
		rt := internal.NewLoggingRoundTripper(&dummyTransport{}, &buf)

		req, err := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		assert.NoError(t, err)
		req.Header.Set("Hsdp-Api-Signature", "super-secret-signature-value")

		_, err = rt.RoundTrip(req)
		assert.NoError(t, err)

		logOutput := buf.String()
		assert.NotContains(t, logOutput, "super-secret-signature-value")
		assert.Contains(t, logOutput, "Hsdp-Api-Signature: [sensitive]")
	})

	t.Run("Hsdp-Api-Signature lowercase request header redaction", func(t *testing.T) {
		var buf bytes.Buffer
		rt := internal.NewLoggingRoundTripper(&dummyTransport{}, &buf)

		req, err := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		assert.NoError(t, err)
		req.Header["hsdp-api-signature"] = []string{"lowercase-signature-value"}

		_, err = rt.RoundTrip(req)
		assert.NoError(t, err)

		logOutput := buf.String()
		assert.NotContains(t, logOutput, "lowercase-signature-value")
		assert.Contains(t, logOutput, "Hsdp-Api-Signature: [sensitive]")
	})

	t.Run("Hsdp-Api-Signature response header redaction", func(t *testing.T) {
		var buf bytes.Buffer
		transport := &dummyTransport{
			responseHeader: http.Header{
				"Hsdp-Api-Signature": []string{"response-secret-signature"},
			},
		}
		rt := internal.NewLoggingRoundTripper(transport, &buf)

		req, err := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		assert.NoError(t, err)

		_, err = rt.RoundTrip(req)
		assert.NoError(t, err)

		logOutput := buf.String()
		assert.NotContains(t, logOutput, "response-secret-signature")
		assert.Contains(t, logOutput, "Hsdp-Api-Signature: [sensitive]")
	})

	t.Run("Form-encoded refresh_token redaction", func(t *testing.T) {
		var buf bytes.Buffer
		rt := internal.NewLoggingRoundTripper(&dummyTransport{}, &buf)

		body := "grant_type=refresh_token&refresh_token=secret_refresh-token.value%2012345&scope=read"
		req, err := http.NewRequest(http.MethodPost, "http://example.com/oauth/token", strings.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		_, err = rt.RoundTrip(req)
		assert.NoError(t, err)

		logOutput := buf.String()
		assert.NotContains(t, logOutput, "secret_refresh-token.value%2012345")
		assert.NotContains(t, logOutput, "-token.value%2012345")
		assert.Contains(t, logOutput, "refresh_token=sensitive")
	})

	t.Run("JSON refresh_token redaction", func(t *testing.T) {
		var buf bytes.Buffer
		rt := internal.NewLoggingRoundTripper(&dummyTransport{}, &buf)

		body := `{"grant_type":"refresh_token","refresh_token":"json_secret_refresh_token"}`
		req, err := http.NewRequest(http.MethodPost, "http://example.com/oauth/token", strings.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		_, err = rt.RoundTrip(req)
		assert.NoError(t, err)

		logOutput := buf.String()
		assert.NotContains(t, logOutput, "json_secret_refresh_token")
		assert.Contains(t, logOutput, `"refresh_token":"[sensitive]"`)
	})
}
