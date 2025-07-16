package model

// Task represents a task item
type Task struct {
	ID     string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name   string `json:"name" example:"Learn Go programming"`
	Status int    `json:"status" example:"0" enums:"0,1"`
}

// TaskRequest represents the request payload for creating or updating a task
type TaskRequest struct {
	Name   string `json:"name" binding:"required" example:"Learn Go programming"`
	Status int    `json:"status" binding:"required" example:"0" enums:"0,1"`
}

// ErrorResponse represents error response format
type ErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}

// MessageResponse represents success message response format
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}