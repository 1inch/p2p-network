// Package types includes helper type definitions
package types

// JsonRequest describes payload for JSON-RPC request
type JsonRequest struct {
	Id     string   `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"` // TODO maybe should change this to interface{} to pass parameters with field names
}

// JsonResponse describes payload for JSON-RPC response
type JsonResponse struct {
	Id     string      `json:"id"`
	Result interface{} `json:"result"`
}
