package task

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/model"
	"github.com/gogolook/task-api/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListTasks(t *testing.T) {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockStorage    *storage.MockStorage
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "成功取得空清單",
			mockStorage: &storage.MockStorage{
				ListFunc: func() []model.Task {
					return []model.Task{}
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`,
		},
		{
			name: "成功取得資料清單",
			mockStorage: &storage.MockStorage{
				ListFunc: func() []model.Task {
					return []model.Task{
						{ID: "1", Name: "Task 1", Status: 0},
						{ID: "2", Name: "Task 2", Status: 1},
					}
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":"1","name":"Task 1","status":0},{"id":"2","name":"Task 2","status":1}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 建立 handler
			handler := NewTaskHandler(tt.mockStorage)

			// 建立 request
			req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
			require.NoError(t, err)

			// 建立 response recorder
			w := httptest.NewRecorder()

			// 建立 gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 執行 handler
			handler.ListTasks(c)

			// 檢查 status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 檢查 response body
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}