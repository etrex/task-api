package task

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/model"
	"github.com/gogolook/task-api/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		taskID         string
		setupTask      *model.Task
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "成功取得任務",
			taskID: "test-id",
			setupTask: &model.Task{
				ID:     "test-id",
				Name:   "Test Task",
				Status: 0,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"test-id","name":"Test Task","status":0}`,
		},
		{
			name:           "任務不存在",
			taskID:         "non-existent",
			setupTask:      nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Task not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 設定
			storage := storage.NewMemoryStorage()
			handler := NewTaskHandler(storage)

			// 如果有設定任務，先創建它
			if tt.setupTask != nil {
				storage.Create(tt.setupTask)
			}

			router := gin.New()
			router.GET("/tasks/:id", handler.GetTask)

			// 執行
			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tt.taskID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 驗證
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}