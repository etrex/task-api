package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/model"
)

// CreateTask 處理建立新資料的 HTTP 請求
// @Summary Create a new task
// @Description Create a new task with name and status
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body model.TaskRequest true "Task data"
// @Success 201 {object} model.Task
// @Failure 400 {object} model.BadRequestResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task model.Task
	
	// 驗證請求資料
	if err := validateTaskRequest(c, &task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 嘗試寫入到 storage，若失敗回傳伺服器錯誤
	if err := h.storage.Create(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	// 回傳建立成功的資料
	c.JSON(http.StatusCreated, task)
}