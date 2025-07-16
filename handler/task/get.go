package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/storage"
)

// GetTask 處理取得單一任務的 HTTP 請求
// @Summary Get a task by ID
// @Description Get a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} model.Task
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /tasks/{id} [get]
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