package main

import (
	"net/http"
	
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
	r.DELETE("/tasks", taskHandler.DeleteAllTasks)
	
	// 健康檢查 endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}