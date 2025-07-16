package task

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/model"
)

// validateTaskRequest 驗證 Task 請求
func validateTaskRequest(c *gin.Context, task *model.Task) error {
	// 先解析到 raw map 檢查必填欄位是否存在
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// 檢查必填欄位是否存在
	if _, exists := raw["name"]; !exists {
		return errors.New("name is required")
	}

	if _, exists := raw["status"]; !exists {
		return errors.New("status is required")
	}

	// 檢查型別並賦值
	name, ok := raw["name"].(string)
	if !ok {
		return errors.New("name must be a string")
	}
	
	// 檢查 name 不能為空字串或僅包含空白字元
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}
	
	task.Name = name

	status, ok := raw["status"].(float64)
	if !ok {
		return errors.New("status must be a number")
	}
	
	// 檢查 status 範圍
	if status < 0 || status > 1 {
		return errors.New("status must be 0 or 1")
	}
	
	task.Status = int(status)

	return nil
}