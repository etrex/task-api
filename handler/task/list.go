package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/storage"
)

// ListTasks 處理列出所有資料的 HTTP 請求
func (h *TaskHandler) ListTasks(c *gin.Context) {
	// 解析分頁參數（只允許传递 page）
	pageStr := c.DefaultQuery("page", "1")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	// 建立分頁參數（後端固定每頁 100 筆）
	params := storage.NewPaginationParams(page)
	
	result, err := h.storage.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}