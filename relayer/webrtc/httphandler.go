package webrtc

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pion/webrtc/v4"
)

// SDPHandler handles SDP request.
func SDPHandler(logger *slog.Logger, sdpRequests chan SDPRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			SessionID string                    `json:"session_id"`
			Offer     webrtc.SessionDescription `json:"offer"`
		}
		logger.Info("sdp request")

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		responseChan := make(chan *webrtc.SessionDescription)
		sdpRequests <- SDPRequest{
			SessionID: req.SessionID,
			Offer:     req.Offer,
			Response:  responseChan,
		}

		answer := <-responseChan
		if answer == nil {
			http.Error(w, "failed to process sdp offer", http.StatusInternalServerError)
			return
		}

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
			SessionID string                  `json:"session_id"`
			Candidate webrtc.ICECandidateInit `json:"candidate"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("candidate parse error", "err", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		log.Info("### CandidateHandler", "sessionID", req.SessionID)
		log.Info("### CandidateHandler", "candidate", req.Candidate)

		candidates <- ICECandidate{
			SessionID: req.SessionID,
			Candidate: req.Candidate,
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusAccepted)
	}
}
