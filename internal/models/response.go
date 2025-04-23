package models

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"` // Use `omitempty` to exclude if not provided
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
