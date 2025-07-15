package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/handler/task"
	"github.com/gogolook/task-api/storage"
)

func main() {
	r := gin.Default()

	memStorage := storage.NewMemoryStorage()
	taskHandler := task.NewTaskHandler(memStorage)

	r.GET("/tasks", taskHandler.ListTasks)
	r.POST("/tasks", taskHandler.CreateTask)
	r.PUT("/tasks/:id", taskHandler.UpdateTask)
	r.DELETE("/tasks/:id", taskHandler.DeleteTask)

	r.Run(":8080")
}