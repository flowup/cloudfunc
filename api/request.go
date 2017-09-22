package api

import "encoding/json"

// Request is structure that represents http request that comes to function
type Request struct {
	BaseURL string `json:"baseUrl"`
	Body json.RawMessage `json:"body"`
	Hostname string `json:"hostname"`
	IP string `json:"ip"`
	Method string `json:"method"`
	OriginalURL string `json:"originalUrl"`
	Query map[string]string `json:"query"`
}

// BindBody binds request's body to given interface
func (r *Request) BindBody(i interface{}) error {
	err := json.Unmarshal(r.Body, i)
	return err
}