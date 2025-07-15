package task

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteTask(t *testing.T) {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		taskID         string
		mockStorage    *storage.MockStorage
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "成功刪除資料",
			taskID: "test-id-123",
			mockStorage: &storage.MockStorage{
				DeleteFunc: func(id string) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"task deleted successfully"}`,
		},
		{
			name:   "資料不存在",
			taskID: "non-existing-id",
			mockStorage: &storage.MockStorage{
				DeleteFunc: func(id string) error {
					return storage.ErrTaskNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"task not found"}`,
		},
		{
			name:   "Storage 錯誤",
			taskID: "test-id-123",
			mockStorage: &storage.MockStorage{
				DeleteFunc: func(id string) error {
					return errors.New("storage error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to delete task"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 建立 handler
			handler := NewTaskHandler(tt.mockStorage)

			// 建立 request
			req, err := http.NewRequest(http.MethodDelete, "/tasks/"+tt.taskID, nil)
			require.NoError(t, err)

			// 建立 response recorder
			w := httptest.NewRecorder()

			// 建立 gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{
				{Key: "id", Value: tt.taskID},
			}

			// 執行 handler
			handler.DeleteTask(c)

			// 檢查 status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 檢查 response body
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}