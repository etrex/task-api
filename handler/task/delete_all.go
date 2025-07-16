package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteAllTasks 處理刪除所有任務的 HTTP 請求
// @Summary Delete all tasks
// @Description Delete all tasks (testing utility)
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} model.MessageResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /tasks [delete]
func (h *TaskHandler) DeleteAllTasks(c *gin.Context) {
	err := h.storage.DeleteAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "All tasks deleted successfully"})
}