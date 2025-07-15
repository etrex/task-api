package model

type Task struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}