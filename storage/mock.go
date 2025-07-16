package storage

import (
	"github.com/gogolook/task-api/model"
)

// MockStorage 用於測試的 mock storage
type MockStorage struct {
	ListFunc   func(params PaginationParams) (*PaginationResult, error)
	GetFunc    func(id string) (*model.Task, error)
	CreateFunc func(task *model.Task) error
	UpdateFunc func(id string, task *model.Task) error
	DeleteFunc func(id string) error
}

func (m *MockStorage) List(params PaginationParams) (*PaginationResult, error) {
	if m.ListFunc != nil {
		return m.ListFunc(params)
	}
	return &PaginationResult{
		Data: []model.Task{},
		Pagination: PaginationInfo{
			Page:    params.Page,
			Limit:   params.Limit,
			Total:   0,
			Pages:   0,
			HasNext: false,
			HasPrev: false,
		},
	}, nil
}

func (m *MockStorage) Get(id string) (*model.Task, error) {
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return nil, nil
}

func (m *MockStorage) Create(task *model.Task) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(task)
	}
	return nil
}

func (m *MockStorage) Update(id string, task *model.Task) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, task)
	}
	return nil
}

func (m *MockStorage) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}