package webrtc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/assert"
)

// TestSDPHandler_Success tests the successful SDP response flow.
func TestSDPHandler_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan SDPRequest, 1)

	go func() {
		for req := range sdpRequests {
			// Simulate processing and sending a valid SDP answer.
			answer := &webrtc.SessionDescription{
				Type: webrtc.SDPTypeAnswer,
				SDP:  "v=0\r\no=- 12345 2 IN IP4 127.0.0.1\r\n",
			}
			req.Response <- answer
		}
	}()

	handler := SDPHandler(logger, sdpRequests)

	// Prepare the HTTP request.
	payload := map[string]interface{}{
		"session_id": "test-session",
		"offer": map[string]string{
			"type": "offer",
			"sdp":  "v=0\r\no=- 54321 2 IN IP4 127.0.0.1\r\n",
		},
	}
	body, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/sdp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler(rec, req)

	// Validate the response.
	assert.Equal(t, http.StatusOK, rec.Code)
	var response struct {
		Answer webrtc.SessionDescription `json:"answer"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, webrtc.SDPTypeAnswer, response.Answer.Type)
	assert.NotEmpty(t, response.Answer.SDP)
}

// TestSDPHandler_InvalidBody tests the case when the request body is invalid.
func TestSDPHandler_InvalidBody(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan SDPRequest, 1)

	handler := SDPHandler(logger, sdpRequests)

	// Invalid body.
	req := httptest.NewRequest(http.MethodPost, "/sdp", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler(rec, req)

	// Validate the response.
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid request body")
}

// TestSDPHandler_FailedResponse tests the case where the server returns a nil response.
func TestSDPHandler_FailedResponse(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan SDPRequest, 1)

	// Simulate the server sending a nil response.
	go func() {
		for req := range sdpRequests {
			req.Response <- nil
		}
	}()

	handler := SDPHandler(logger, sdpRequests)

	// Prepare the HTTP request.
	payload := map[string]interface{}{
		"session_id": "test-session",
		"offer": map[string]string{
			"type": "offer",
			"sdp":  "v=0\r\no=- 54321 2 IN IP4 127.0.0.1\r\n",
		},
	}
	body, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/sdp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler(rec, req)

	// Validate the response.
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to process sdp offer")
}
