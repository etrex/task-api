package model

type Task struct {
	ID     string `json:"id"`
	Name   string `json:"name" binding:"required"`
	Status int    `json:"status" binding:"min=0,max=1"`
}