package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteAllTasks 處理刪除所有任務的 HTTP 請求
func (h *TaskHandler) DeleteAllTasks(c *gin.Context) {
	err := h.storage.DeleteAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "All tasks deleted successfully"})
}