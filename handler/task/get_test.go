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
			taskID: "", // 將在測試中設定
			setupTask: &model.Task{
				Name:   "Test Task",
				Status: 0,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "", // 將在測試中設定
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

			var actualTaskID string
			var expectedBody string

			// 如果有設定任務，先創建它
			if tt.setupTask != nil {
				storage.Create(tt.setupTask)
				// 獲取實際生成的 ID
				tasks := storage.List()
				if len(tasks) > 0 {
					actualTaskID = tasks[0].ID
					expectedBody = `{"id":"` + actualTaskID + `","name":"Test Task","status":0}`
				}
			} else {
				actualTaskID = tt.taskID
				expectedBody = tt.expectedBody
			}

			router := gin.New()
			router.GET("/tasks/:id", handler.GetTask)

			// 執行
			req := httptest.NewRequest(http.MethodGet, "/tasks/"+actualTaskID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 驗證
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, expectedBody, w.Body.String())
		})
	}
}