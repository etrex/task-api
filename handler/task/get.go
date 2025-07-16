package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/storage"
)

func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	
	task, err := h.storage.Get(id)
	if err != nil {
		if err == storage.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	
	c.JSON(http.StatusOK, task)
}