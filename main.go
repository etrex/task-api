package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/handler/task"
	"github.com/gogolook/task-api/storage"
)

func main() {
	r := gin.Default()

	// CORS middleware - 只允許前端域名
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "https://etrex.tw" {
			c.Header("Access-Control-Allow-Origin", "https://etrex.tw")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")
		}
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	memStorage := storage.NewMemoryStorage()
	taskHandler := task.NewTaskHandler(memStorage)

	r.GET("/tasks", taskHandler.ListTasks)
	r.GET("/tasks/:id", taskHandler.GetTask)
	r.POST("/tasks", taskHandler.CreateTask)
	r.PUT("/tasks/:id", taskHandler.UpdateTask)
	r.DELETE("/tasks/:id", taskHandler.DeleteTask)

	r.Run(":8080")
}