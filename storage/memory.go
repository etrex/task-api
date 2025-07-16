package storage

import (
	"errors"
	"sync"

	"github.com/gogolook/task-api/model"
	"github.com/google/uuid"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

// PaginationParams 分頁參數
type PaginationParams struct {
	Page  int
	Limit int
}

// NewPaginationParams 建立分頁參數（限制每頁 100 筆）
func NewPaginationParams(page int) PaginationParams {
	if page < 1 {
		page = 1
	}
	return PaginationParams{
		Page:  page,
		Limit: 100,
	}
}

// PaginationResult 分頁結果
type PaginationResult struct {
	Data   []model.Task `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo 分頁資訊
type PaginationInfo struct {
	Page     int  `json:"page"`
	Limit    int  `json:"limit"`
	Total    int  `json:"total"`
	Pages    int  `json:"pages"`
	HasNext  bool `json:"has_next"`
	HasPrev  bool `json:"has_prev"`
}

type Storage interface {
	List(params PaginationParams) (*PaginationResult, error)
	Get(id string) (*model.Task, error)
	Create(task *model.Task) error
	Update(id string, task *model.Task) error
	Delete(id string) error
}

type MemoryStorage struct {
	mu        sync.RWMutex
	tasks     []model.Task      // 使用 slice 儲存，保持插入順序
	indexMap  map[string]int    // uuid -> slice index 的映射
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks:    make([]model.Task, 0),
		indexMap: make(map[string]int),
	}
}

func (s *MemoryStorage) List(params PaginationParams) (*PaginationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// 驗證參數
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 100
	}
	
	total := len(s.tasks)
	pages := (total + params.Limit - 1) / params.Limit // 向上取整
	
	// 計算 offset
	offset := (params.Page - 1) * params.Limit
	
	// 直接從 slice 取得分頁資料 - O(limit)
	var data []model.Task
	if offset < total {
		end := offset + params.Limit
		if end > total {
			end = total
		}
		// 直接切片，O(1) 操作
		data = make([]model.Task, end-offset)
		copy(data, s.tasks[offset:end])
	} else {
		data = []model.Task{}
	}
	
	// 建立分頁資訊
	pagination := PaginationInfo{
		Page:    params.Page,
		Limit:   params.Limit,
		Total:   total,
		Pages:   pages,
		HasNext: params.Page < pages,
		HasPrev: params.Page > 1,
	}
	
	return &PaginationResult{
		Data:       data,
		Pagination: pagination,
	}, nil
}

func (s *MemoryStorage) Get(id string) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	index, exists := s.indexMap[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	
	task := s.tasks[index]
	return &task, nil
}

func (s *MemoryStorage) Create(task *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	task.ID = uuid.New().String()
	
	// 新增到 slice 的最後
	s.tasks = append(s.tasks, *task)
	
	// 更新 index map
	s.indexMap[task.ID] = len(s.tasks) - 1
	
	return nil
}

func (s *MemoryStorage) Update(id string, task *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	index, exists := s.indexMap[id]
	if !exists {
		return ErrTaskNotFound
	}

	task.ID = id
	s.tasks[index] = *task
	
	return nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	index, exists := s.indexMap[id]
	if !exists {
		return ErrTaskNotFound
	}
	
	lastIndex := len(s.tasks) - 1
	
	// 如果不是最後一個元素，將最後一個元素移到被刪除的位置
	if index != lastIndex {
		s.tasks[index] = s.tasks[lastIndex]
		// 更新被移動元素的 index
		s.indexMap[s.tasks[index].ID] = index
	}
	
	// 刪除最後一個元素
	s.tasks = s.tasks[:lastIndex]
	
	// 從 index map 中刪除
	delete(s.indexMap, id)
	
	return nil
}