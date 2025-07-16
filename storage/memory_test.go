package storage

import (
	"testing"

	"github.com/gogolook/task-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_Create(t *testing.T) {
	storage := NewMemoryStorage()
	
	task := &model.Task{
		Name:   "Test Task",
		Status: 0,
	}
	
	err := storage.Create(task)
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)
}

func TestMemoryStorage_List(t *testing.T) {
	storage := NewMemoryStorage()
	
	// 新增測試資料
	tasks := []*model.Task{
		{Name: "Task 1", Status: 0},
		{Name: "Task 2", Status: 1},
		{Name: "Task 3", Status: 0},
	}
	
	for _, task := range tasks {
		err := storage.Create(task)
		require.NoError(t, err)
	}

	// 測試 List
	result, err := storage.List(NewPaginationParams(1))
	require.NoError(t, err)
	assert.Len(t, result.Data, 3)
	assert.Equal(t, 3, result.Pagination.Total)
	assert.Equal(t, 1, result.Pagination.Pages)
	assert.False(t, result.Pagination.HasNext)
	assert.False(t, result.Pagination.HasPrev)
}

func TestMemoryStorage_Get(t *testing.T) {
	storage := NewMemoryStorage()
	
	// 新增測試資料
	task := &model.Task{
		Name:   "Test Task",
		Status: 0,
	}
	
	err := storage.Create(task)
	require.NoError(t, err)
	
	// 測試 Get
	retrieved, err := storage.Get(task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.Name, retrieved.Name)
	assert.Equal(t, task.Status, retrieved.Status)
	assert.Equal(t, task.ID, retrieved.ID)
}

func TestMemoryStorage_Update(t *testing.T) {
	storage := NewMemoryStorage()
	
	// 新增測試資料
	task := &model.Task{
		Name:   "Test Task",
		Status: 0,
	}
	
	err := storage.Create(task)
	require.NoError(t, err)
	
	// 測試 Update
	updatedTask := &model.Task{
		Name:   "Updated Task",
		Status: 1,
	}
	
	err = storage.Update(task.ID, updatedTask)
	require.NoError(t, err)
	
	// 驗證更新結果
	retrieved, err := storage.Get(task.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Task", retrieved.Name)
	assert.Equal(t, 1, retrieved.Status)
	assert.Equal(t, task.ID, retrieved.ID)
}

func TestMemoryStorage_Delete(t *testing.T) {
	storage := NewMemoryStorage()
	
	// 新增測試資料
	task := &model.Task{
		Name:   "Test Task",
		Status: 0,
	}
	
	err := storage.Create(task)
	require.NoError(t, err)
	
	// 測試 Delete
	err = storage.Delete(task.ID)
	require.NoError(t, err)
	
	// 驗證刪除結果
	_, err = storage.Get(task.ID)
	assert.Equal(t, ErrTaskNotFound, err)
}

func TestMemoryStorage_DeleteAll(t *testing.T) {
	storage := NewMemoryStorage()
	
	// 新增測試資料
	tasks := []*model.Task{
		{Name: "Task 1", Status: 0},
		{Name: "Task 2", Status: 1},
		{Name: "Task 3", Status: 0},
	}
	
	for _, task := range tasks {
		err := storage.Create(task)
		require.NoError(t, err)
	}
	
	// 驗證任務已新增
	result, err := storage.List(NewPaginationParams(1))
	require.NoError(t, err)
	assert.Equal(t, 3, result.Pagination.Total)
	
	// 測試 DeleteAll
	err = storage.DeleteAll()
	require.NoError(t, err)
	
	// 驗證所有任務已刪除
	result, err = storage.List(NewPaginationParams(1))
	require.NoError(t, err)
	assert.Equal(t, 0, result.Pagination.Total)
	assert.Len(t, result.Data, 0)
}

func TestMemoryStorage_GetNotFound(t *testing.T) {
	storage := NewMemoryStorage()
	
	_, err := storage.Get("nonexistent")
	assert.Equal(t, ErrTaskNotFound, err)
}

func TestMemoryStorage_UpdateNotFound(t *testing.T) {
	storage := NewMemoryStorage()
	
	task := &model.Task{
		Name:   "Test Task",
		Status: 0,
	}
	
	err := storage.Update("nonexistent", task)
	assert.Equal(t, ErrTaskNotFound, err)
}

func TestMemoryStorage_DeleteNotFound(t *testing.T) {
	storage := NewMemoryStorage()
	
	err := storage.Delete("nonexistent")
	assert.Equal(t, ErrTaskNotFound, err)
}