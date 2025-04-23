package webrtc

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/1inch/p2p-network/relayer/metrics"
	"github.com/pion/webrtc/v4"
)

// SDPHandler handles SDP request.
func SDPHandler(logger *slog.Logger, sdpRequests chan SDPRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		origin := r.Header.Get("Origin")
		if origin == "" {
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			origin = scheme + "://" + r.Host
		}
		candidateURL := origin + "/candidate"

		var req struct {
			SessionID string                    `json:"session_id"`
			Offer     webrtc.SessionDescription `json:"offer"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		responseChan := make(chan *webrtc.SessionDescription)
		sdpRequests <- SDPRequest{
			SessionID:    req.SessionID,
			Offer:        req.Offer,
			CandidateURL: candidateURL,
			Response:     responseChan,
		}

		answer := <-responseChan
		if answer == nil {
			http.Error(w, "failed to process sdp offer", http.StatusInternalServerError)
			return
		}

		latency := time.Since(start).Seconds()
		metrics.EndToEndWorkflowLatency.Observe(latency)
		metrics.EndToEndWorkflowCompleted.Inc()

		resp := struct {
			Answer webrtc.SessionDescription `json:"answer"`
		}{Answer: *answer}

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// CandidateHandler handles ICECandidate request.
func CandidateHandler(log *slog.Logger, candidates chan ICECandidate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			SessionID string              `json:"session_id"`
			Candidate webrtc.ICECandidate `json:"candidate"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		candidates <- ICECandidate{
			SessionID: req.SessionID,
			Candidate: req.Candidate,
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
