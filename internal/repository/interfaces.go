package repository

import (
	"context"
	"github.com/SaveljevRoman/go-layout-project-2/internal/domain/entity"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id string) (*entity.Task, error)
	GetAll(ctx context.Context, userID string) ([]*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id string) error
}

type LogRepository interface {
	LogInfo(message string, fields map[string]interface{})
	LogError(message string, err error, fields map[string]interface{})
}
