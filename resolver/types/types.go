package types

type JsonRequest struct {
	Id     string   `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type JsonResponse struct {
	Id     string      `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}
