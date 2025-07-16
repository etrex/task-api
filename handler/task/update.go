package task

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/model"
	"github.com/gogolook/task-api/storage"
)

// UpdateTask 處理更新指定資料的 HTTP 請求
// @Summary Update a task
// @Description Update a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body model.TaskRequest true "Task data"
// @Success 200 {object} model.Task
// @Failure 400 {object} model.BadRequestResponse
// @Failure 404 {object} model.NotFoundResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	
	var task model.Task
	// 驗證請求資料
	if err := validateTaskRequest(c, &task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新 storage 中的資料，若資料不存在回傳 404，其他錯誤回傳 500
	if err := h.storage.Update(id, &task); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	// 回傳更新後的資料
	c.JSON(http.StatusOK, task)
}