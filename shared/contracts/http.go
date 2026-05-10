package contracts

// APIResponse is the response structure for the API.
type APIResponse struct {
	Ok      bool      `json:"ok"`
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

// APIError is the error structure for the API.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
