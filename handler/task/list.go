package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListTasks 處理列出所有資料的 HTTP 請求
func (h *TaskHandler) ListTasks(c *gin.Context) {
	// 從 storage 取得所有資料並回傳 JSON 格式的清單
	tasks := h.storage.List()
	c.JSON(http.StatusOK, tasks)
}