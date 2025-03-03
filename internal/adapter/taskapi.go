package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SaveljevRoman/go-layout-project-2/internal/domain/entity"
	"github.com/SaveljevRoman/go-layout-project-2/internal/usecase"
)

type TaskAPI struct {
	taskUseCase *usecase.TaskUseCase
}

func NewTaskAPI(taskUseCase *usecase.TaskUseCase) *TaskAPI {
	return &TaskAPI{
		taskUseCase: taskUseCase,
	}
}

// SyncWithExternalAPI Пример метода для интеграции с внешним API
func (a *TaskAPI) SyncWithExternalAPI(ctx context.Context, externalTasks []entity.Task) error {
	// Это просто пример адаптера, который может синхронизировать задачи с внешним API
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return errors.New("unauthorized")
	}

	for _, extTask := range externalTasks {
		task := entity.Task{
			Title:       extTask.Title,
			Description: extTask.Description,
			Status:      extTask.Status,
		}

		if err := a.taskUseCase.CreateTask(ctx, &task); err != nil {
			return err
		}
	}

	return nil
}

// ExportTasksToJSON Метод для экспорта задач в JSON
func (a *TaskAPI) ExportTasksToJSON(ctx context.Context) ([]byte, error) {
	tasks, err := a.taskUseCase.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	return json.Marshal(tasks)
}
