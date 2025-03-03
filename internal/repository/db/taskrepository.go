package db

import (
	"context"
	"errors"
	"github.com/SaveljevRoman/go-layout-project-2/internal/domain/entity"
	"sync"
	"time"
)

type TaskRepository struct {
	// В реальном приложении здесь будет подключение к БД
	tasks map[string]*entity.Task
	mutex sync.RWMutex
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*entity.Task),
	}
}

func (r *TaskRepository) Create(ctx context.Context, task *entity.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Имитация генерации ID
	if task.ID == "" {
		task.ID = time.Now().Format("20060102150405")
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*entity.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	return task, nil
}

func (r *TaskRepository) GetAll(ctx context.Context, userID string) ([]*entity.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*entity.Task

	for _, task := range r.tasks {
		if task.UserID == userID {
			result = append(result, task)
		}
	}

	return result, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *entity.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.tasks[task.ID]
	if !exists {
		return errors.New("task not found")
	}

	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	delete(r.tasks, id)

	return nil
}
