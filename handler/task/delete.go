package task

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/storage"
)

// DeleteTask 處理刪除指定資料的 HTTP 請求
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	// 從 storage 刪除資料，若資料不存在回傳 404，其他錯誤回傳 500
	if err := h.storage.Delete(id); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	// 回傳刪除成功訊息
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}