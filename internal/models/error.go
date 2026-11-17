package models

type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Details any    `json:"details,omitempty"`
}
