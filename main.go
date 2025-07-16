package main

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/handler/task"
	"github.com/gogolook/task-api/storage"
)

// @title Task API
// @version 1.0
// @description A RESTful task management API with high-performance in-memory storage
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host task-api.etrex.tw
// @BasePath /
// @schemes https http

func main() {
	r := gin.Default()

	// CORS middleware - 允許前端域名和 GitHub Pages
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "https://etrex.tw" || origin == "https://etrex.github.io" {
			c.Header("Access-Control-Allow-Origin", origin)
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