package storage

import (
	"testing"

	"github.com/gogolook/task-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_Create(t *testing.T) {
	// 初始化 storage
	storage := NewMemoryStorage()

	// 測試建立資料
	task := &model.Task{
		Name:   "Test Task",
		Status: 0,
	}

	err := storage.Create(task)
	require.NoError(t, err)

	// 確認 ID 有被產生
	assert.NotEmpty(t, task.ID)

	// 確認資料有被儲存
	storedTask, err := storage.Get(task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.Name, storedTask.Name)
}

func TestMemoryStorage_List(t *testing.T) {
	// 初始化 storage
	storage := NewMemoryStorage()

	// 建立多筆資料
	tasks := []model.Task{
		{Name: "Task 1", Status: 0},
		{Name: "Task 2", Status: 1},
		{Name: "Task 3", Status: 0},
	}

	for i := range tasks {
		storage.Create(&tasks[i])
	}

	// 測試 List
	result := storage.List()
	assert.Len(t, result, 3)
}

func TestMemoryStorage_Get(t *testing.T) {
	// 初始化 storage
	storage := NewMemoryStorage()

	// 建立測試資料
	task := &model.Task{
		Name:   "Test Task",
		Status: 1,
	}
	storage.Create(task)

	// 測試成功取得資料
	t.Run("existing task", func(t *testing.T) {
		result, err := storage.Get(task.ID)
		require.NoError(t, err)
		assert.Equal(t, task.ID, result.ID)
	})

	// 測試取得不存在的資料
	t.Run("non-existing task", func(t *testing.T) {
		_, err := storage.Get("non-existing-id")
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}

func TestMemoryStorage_Update(t *testing.T) {
	// 初始化 storage
	storage := NewMemoryStorage()

	// 建立測試資料
	task := &model.Task{
		Name:   "Original Task",
		Status: 0,
	}
	storage.Create(task)
	originalID := task.ID

	// 測試成功更新
	t.Run("existing task", func(t *testing.T) {
		updatedTask := &model.Task{
			Name:   "Updated Task",
			Status: 1,
		}

		err := storage.Update(originalID, updatedTask)
		require.NoError(t, err)

		// 確認資料有被更新
		result, _ := storage.Get(originalID)
		assert.Equal(t, "Updated Task", result.Name)
		assert.Equal(t, 1, result.Status)
		assert.Equal(t, originalID, result.ID)
	})

	// 測試更新不存在的資料
	t.Run("non-existing task", func(t *testing.T) {
		err := storage.Update("non-existing-id", &model.Task{})
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}

func TestMemoryStorage_Delete(t *testing.T) {
	// 初始化 storage
	storage := NewMemoryStorage()

	// 建立測試資料
	task := &model.Task{
		Name:   "Task to Delete",
		Status: 0,
	}
	storage.Create(task)

	// 測試成功刪除
	t.Run("existing task", func(t *testing.T) {
		err := storage.Delete(task.ID)
		require.NoError(t, err)

		// 確認資料已被刪除
		_, err = storage.Get(task.ID)
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})

	// 測試刪除不存在的資料
	t.Run("non-existing task", func(t *testing.T) {
		err := storage.Delete("non-existing-id")
		assert.ErrorIs(t, err, ErrTaskNotFound)
	})
}