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

func TestCreateTask(t *testing.T) {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockStorage    *storage.MockStorage
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "成功建立資料",
			requestBody: map[string]interface{}{
				"name":   "Test Task",
				"status": 0,
			},
			mockStorage: &storage.MockStorage{
				CreateFunc: func(task *model.Task) error {
					task.ID = "test-id-123"
					return nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"test-id-123","name":"Test Task","status":0}`,
		},
		{
			name: "JSON 解析錯誤",
			requestBody: `{invalid json}`,
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name: "缺少 name 欄位",
			requestBody: map[string]interface{}{
				"status": 0,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name: "缺少 status 欄位",
			requestBody: map[string]interface{}{
				"name": "Test Task",
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name: "status 值超出範圍 (小於 0)",
			requestBody: map[string]interface{}{
				"name":   "Test Task",
				"status": -1,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name: "status 值超出範圍 (大於 1)",
			requestBody: map[string]interface{}{
				"name":   "Test Task",
				"status": 2,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name: "status 型別錯誤 (字串)",
			requestBody: map[string]interface{}{
				"name":   "Test Task",
				"status": "1",
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name: "name 型別錯誤 (數字)",
			requestBody: map[string]interface{}{
				"name":   123,
				"status": 1,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name: "name 為空字串",
			requestBody: map[string]interface{}{
				"name":   "",
				"status": 0,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name: "name 為純空白字串",
			requestBody: map[string]interface{}{
				"name":   "   ",
				"status": 0,
			},
			mockStorage: &storage.MockStorage{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":`,
		},
		{
			name: "Storage 錯誤",
			requestBody: map[string]interface{}{
				"name":   "Test Task",
				"status": 0,
			},
			mockStorage: &storage.MockStorage{
				CreateFunc: func(task *model.Task) error {
					return errors.New("storage error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to create task"}`,
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
			req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// 建立 response recorder
			w := httptest.NewRecorder()

			// 建立 gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 執行 handler
			handler.CreateTask(c)

			// 檢查 status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 檢查 response body
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}