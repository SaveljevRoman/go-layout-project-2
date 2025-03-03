package entity

import (
	"time"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "TODO"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	UserID      string     `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Validate Валидация задачи
func (t *Task) Validate() []string {
	var errors []string

	if t.Title == "" {
		errors = append(errors, "title is required")
	}

	if len(t.Title) > 100 {
		errors = append(errors, "title must be less than 100 characters")
	}

	if len(t.Description) > 1000 {
		errors = append(errors, "description must be less than 1000 characters")
	}

	// Проверка статуса
	if t.Status != StatusTodo && t.Status != StatusInProgress && t.Status != StatusDone {
		errors = append(errors, "invalid status")
	}

	return errors
}
