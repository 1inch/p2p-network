// Package metrics provides Prometheus metrics for the relayer
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HttpRequestsTotal Total number of HTTP requests received"
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_http_requests_total",
			Help: "Total number of HTTP requests received",
		},
		[]string{"handler", "method", "status"},
	)
	// HttpRequestDuration Duration of HTTP requests in seconds
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "relayer_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
		[]string{"handler", "method"},
	)

	// IceCandidateSentTotal Total number of ICE candidates sent
	IceCandidateSentTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_ice_candidate_sent_total",
			Help: "Total number of ICE candidates sent",
		},
		[]string{"session_id", "status"},
	)
	// IceCandidateSendDuration Duration for sending ICE candidates in seconds
	IceCandidateSendDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "relayer_ice_candidate_send_duration_seconds",
			Help:    "Duration for sending ICE candidates in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
		[]string{"session_id"},
	)

	// ActivePeerConnections Current number of active PeerConnections
	ActivePeerConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "relayer_active_peer_connections",
			Help: "Current number of active PeerConnections",
		},
	)

	// SdpNegotiationTotal Total number of SDP negotiations
	SdpNegotiationTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_sdp_negotiation_total",
			Help: "Total number of SDP negotiations",
		},
		[]string{"status"},
	)
	// SdpNegotiationDuration Duration of SDP negotiations in seconds
	SdpNegotiationDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "relayer_sdp_negotiation_duration_seconds",
			Help:    "Duration of SDP negotiations in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
	)

	// GrpcRequestsTotal Total number of gRPC requests
	GrpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)
	// GrpcRequestDuration Duration of gRPC requests in seconds
	GrpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "relayer_grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
		[]string{"method"},
	)

	// DataChannelMessagesSent Total number of messages sent over data channels
	DataChannelMessagesSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_data_channel_messages_sent_total",
			Help: "Total number of messages sent over data channels",
		},
		[]string{"session_id", "status"},
	)
	// DataChannelMessagesReceived Total number of messages received over data channels
	DataChannelMessagesReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relayer_data_channel_messages_received_total",
			Help: "Total number of messages received over data channels",
		},
		[]string{"session_id"},
	)
	// DataChannelLatency Duration of data channel message processing in seconds
	DataChannelLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "relayer_data_channel_latency_seconds",
			Help:    "Time taken to process data channel messages in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
		[]string{"session_id"},
	)

	// EndToEndWorkflowLatency Duration of end-to-end workflow in seconds
	EndToEndWorkflowLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "relayer_end_to_end_workflow_latency_seconds",
			Help:    "Latency of end-to-end workflow from receiving SDP offer to sending answer, in seconds",
			Buckets: []float64{0.05, 0.5, 1, 2, 3, 5, 10, 15, 30, 60},
		},
	)
	// EndToEndWorkflowCompleted Total number of completed end-to-end workflows
	EndToEndWorkflowCompleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "relayer_end_to_end_workflow_completed_total",
			Help: "Total number of completed end-to-end workflows",
		},
	)
)

func init() {
	prometheus.MustRegister(
		HttpRequestsTotal, HttpRequestDuration,
		IceCandidateSentTotal, IceCandidateSendDuration,
		ActivePeerConnections,
		SdpNegotiationTotal, SdpNegotiationDuration,
		GrpcRequestsTotal, GrpcRequestDuration,
		DataChannelMessagesSent, DataChannelMessagesReceived,
		DataChannelLatency, EndToEndWorkflowLatency,
		EndToEndWorkflowCompleted,
	)
}

// Handler returns a handler for Prometheus metrics
func Handler() http.Handler {
	return promhttp.Handler()
}
