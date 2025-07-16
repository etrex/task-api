package task

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteAllTasks(t *testing.T) {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockStorage    *storage.MockStorage
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "成功刪除所有任務",
			mockStorage: &storage.MockStorage{
				DeleteAllFunc: func() error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"All tasks deleted successfully"}`,
		},
		{
			name: "刪除失敗",
			mockStorage: &storage.MockStorage{
				DeleteAllFunc: func() error {
					return storage.ErrTaskNotFound
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"task not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 建立 handler
			handler := NewTaskHandler(tt.mockStorage)

			// 建立 request
			req, err := http.NewRequest(http.MethodDelete, "/tasks", nil)
			require.NoError(t, err)

			// 建立 response recorder
			w := httptest.NewRecorder()

			// 建立 gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 執行 handler
			handler.DeleteAllTasks(c)

			// 檢查 status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 檢查 response body
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}