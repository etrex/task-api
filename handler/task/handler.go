package task

import (
	"github.com/gogolook/task-api/storage"
)

type TaskHandler struct {
	storage storage.Storage
}

func NewTaskHandler(storage storage.Storage) *TaskHandler {
	return &TaskHandler{
		storage: storage,
	}
}