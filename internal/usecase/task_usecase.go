package usecase

import (
	"context"
	"errors"
	"github.com/SaveljevRoman/go-layout-project-2/internal/domain/entity"
	"github.com/SaveljevRoman/go-layout-project-2/internal/repository"
	"github.com/SaveljevRoman/go-layout-project-2/pkg/logger"
)

type TaskUseCase struct {
	repo   repository.TaskRepository
	logger *logger.Logger
}

func NewTaskUseCase(repo repository.TaskRepository, logger *logger.Logger) *TaskUseCase {
	return &TaskUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, task *entity.Task) error {
	uc.logger.Info("Creating task", map[string]interface{}{"title": task.Title})

	if errMsgs := task.Validate(); len(errMsgs) > 0 {
		return errors.New("validation failed: " + errMsgs[0])
	}

	// Получаем ID пользователя из контекста (предполагается, что оно туда добавлено middleware)
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return errors.New("unauthorized")
	}

	task.UserID = userID

	return uc.repo.Create(ctx, task)
}

func (uc *TaskUseCase) GetTask(ctx context.Context, id string) (*entity.Task, error) {
	uc.logger.Info("Getting task", map[string]interface{}{"id": id})

	task, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверка прав доступа
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, errors.New("unauthorized")
	}

	if task.UserID != userID {
		return nil, errors.New("access denied")
	}

	return task, nil
}

func (uc *TaskUseCase) GetAllTasks(ctx context.Context) ([]*entity.Task, error) {
	uc.logger.Info("Getting all tasks", nil)

	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, errors.New("unauthorized")
	}

	return uc.repo.GetAll(ctx, userID)
}

func (uc *TaskUseCase) UpdateTask(ctx context.Context, task *entity.Task) error {
	uc.logger.Info("Updating task", map[string]interface{}{"id": task.ID})

	if errMsgs := task.Validate(); len(errMsgs) > 0 {
		return errors.New("validation failed: " + errMsgs[0])
	}

	// Проверка прав доступа
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return errors.New("unauthorized")
	}

	existingTask, err := uc.repo.GetByID(ctx, task.ID)
	if err != nil {
		return err
	}

	if existingTask.UserID != userID {
		return errors.New("access denied")
	}

	task.UserID = userID // Сохраняем оригинального владельца

	return uc.repo.Update(ctx, task)
}

func (uc *TaskUseCase) DeleteTask(ctx context.Context, id string) error {
	uc.logger.Info("Deleting task", map[string]interface{}{"id": id})

	// Проверка прав доступа
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return errors.New("unauthorized")
	}

	existingTask, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingTask.UserID != userID {
		return errors.New("access denied")
	}

	return uc.repo.Delete(ctx, id)
}
