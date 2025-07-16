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

type Storage interface {
	List() []model.Task
	Get(id string) (*model.Task, error)
	Create(task *model.Task) error
	Update(id string, task *model.Task) error
	Delete(id string) error
}

type MemoryStorage struct {
	mu    sync.RWMutex
	tasks map[string]model.Task
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[string]model.Task),
	}
}

func (s *MemoryStorage) List() []model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	result := make([]model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

func (s *MemoryStorage) Get(id string) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	task, exists := s.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	return &task, nil
}

func (s *MemoryStorage) Create(task *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	task.ID = uuid.New().String()
	s.tasks[task.ID] = *task
	return nil
}

func (s *MemoryStorage) Update(id string, task *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	task.ID = id
	s.tasks[id] = *task
	return nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}