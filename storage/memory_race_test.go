package storage

import (
	"fmt"
	"sync"
	"testing"
	"github.com/gogolook/task-api/model"
)

// 使用 go test -race 執行這個測試會檢測到 race condition
func TestMemoryStorage_RaceCondition(t *testing.T) {
	storage := NewMemoryStorage()
	
	t.Run("並發創建任務", func(t *testing.T) {
		var wg sync.WaitGroup
		taskCount := 100
		
		for i := 0; i < taskCount; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				task := &model.Task{
					Name:   fmt.Sprintf("Task %d", index),
					Status: 0,
				}
				storage.Create(task)
			}(i)
		}
		
		wg.Wait()
		
		// 檢查任務數量
		tasks := storage.List()
		if len(tasks) != taskCount {
			t.Errorf("預期 %d 個任務，實際得到 %d 個", taskCount, len(tasks))
		}
	})
	
	t.Run("並發讀寫同一個 map", func(t *testing.T) {
		// 準備測試資料
		task := &model.Task{Name: "Test", Status: 0}
		storage.Create(task)
		taskID := task.ID
		
		var wg sync.WaitGroup
		
		// 同時進行讀取和更新
		for i := 0; i < 50; i++ {
			wg.Add(2)
			
			// 讀取操作
			go func() {
				defer wg.Done()
				storage.Get(taskID)
			}()
			
			// 更新操作
			go func(index int) {
				defer wg.Done()
				updated := &model.Task{
					Name:   fmt.Sprintf("Updated %d", index),
					Status: 1,
				}
				storage.Update(taskID, updated)
			}(i)
		}
		
		wg.Wait()
	})
	
	t.Run("並發刪除和列表", func(t *testing.T) {
		// 創建多個任務
		taskIDs := make([]string, 10)
		for i := 0; i < 10; i++ {
			task := &model.Task{Name: fmt.Sprintf("Task %d", i), Status: 0}
			storage.Create(task)
			taskIDs[i] = task.ID
		}
		
		var wg sync.WaitGroup
		
		// 一半 goroutine 刪除任務
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				storage.Delete(id)
			}(taskIDs[i])
		}
		
		// 另一半 goroutine 讀取列表
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// 在遍歷 map 時，其他 goroutine 正在修改它
				// 這可能導致 "concurrent map iteration and map write" panic
				_ = storage.List()
			}()
		}
		
		wg.Wait()
	})
}