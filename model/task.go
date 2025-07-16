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

// BadRequestResponse represents 400 Bad Request error response format
type BadRequestResponse struct {
	Error string `json:"error" example:"name is required"`
}

// NotFoundResponse represents 404 Not Found error response format
type NotFoundResponse struct {
	Error string `json:"error" example:"task not found"`
}

// MessageResponse represents success message response format
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}