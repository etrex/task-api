package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gogolook/task-api/model"
	"github.com/gogolook/task-api/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateTask(t *testing.T) {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		taskID         string
		requestBody    interface{}
		mockStorage    *storage.MockStorage
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "成功更新資料",
			taskID: "test-id-123",
			requestBody: map[string]interface{}{
				"name":   "Updated Task",
				"status": 1,
			},
			mockStorage: &storage.MockStorage{
				UpdateFunc: func(id string, task *model.Task) error {
					task.ID = id
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"test-id-123","name":"Updated Task","status":1}`,
		},
		{
			name:   "JSON 解析錯誤",
			taskID: "test-id-123",
			requestBody: `{invalid json}`,
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name:   "缺少 name 欄位",
			taskID: "test-id-123",
			requestBody: map[string]interface{}{
				"status": 1,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name:   "缺少 status 欄位",
			taskID: "test-id-123",
			requestBody: map[string]interface{}{
				"name": "Updated Task",
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name:   "status 值超出範圍 (小於 0)",
			taskID: "test-id-123",
			requestBody: map[string]interface{}{
				"name":   "Updated Task",
				"status": -1,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name:   "status 值超出範圍 (大於 1)",
			taskID: "test-id-123",
			requestBody: map[string]interface{}{
				"name":   "Updated Task",
				"status": 2,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name:   "status 型別錯誤 (字串)",
			taskID: "test-id-123",
			requestBody: map[string]interface{}{
				"name":   "Updated Task",
				"status": "1",
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name:   "name 型別錯誤 (數字)",
			taskID: "test-id-123",
			requestBody: map[string]interface{}{
				"name":   123,
				"status": 1,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name:   "資料不存在",
			taskID: "non-existing-id",
			requestBody: map[string]interface{}{
				"name":   "Updated Task",
				"status": 1,
			},
			mockStorage: &storage.MockStorage{
				UpdateFunc: func(id string, task *model.Task) error {
					return storage.ErrTaskNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"task not found"}`,
		},
		{
			name:   "Storage 錯誤",
			taskID: "test-id-123",
			requestBody: map[string]interface{}{
				"name":   "Updated Task",
				"status": 1,
			},
			mockStorage: &storage.MockStorage{
				UpdateFunc: func(id string, task *model.Task) error {
					return errors.New("storage error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to update task"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 建立 handler
			handler := NewTaskHandler(tt.mockStorage)

			// 建立 request body
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			// 建立 request
			req, err := http.NewRequest(http.MethodPut, "/tasks/"+tt.taskID, bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// 建立 response recorder
			w := httptest.NewRecorder()

			// 建立 gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{
				{Key: "id", Value: tt.taskID},
			}

			// 執行 handler
			handler.UpdateTask(c)

			// 檢查 status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 檢查 response body
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}