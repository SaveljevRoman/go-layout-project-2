package handler

import (
	"encoding/json"
	"github.com/SaveljevRoman/go-layout-project-2/internal/domain/entity"
	"github.com/SaveljevRoman/go-layout-project-2/internal/usecase"
	"net/http"
	"strings"
)

type TaskHandler struct {
	taskUseCase *usecase.TaskUseCase
}

func NewTaskHandler(taskUseCase *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
	}
}

// CreateTaskRequest Структуры запросов и ответов
type CreateTaskRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      entity.TaskStatus `json:"status"`
}

type TaskResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      entity.TaskStatus `json:"status"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// CreateTask обрабатывает запрос на создание задачи
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest

	// Валидация запроса
	validationErrors, err := ValidateRequest(r, &req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Если есть ошибки валидации, возвращаем их
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationErrors{Errors: validationErrors})
		return
	}

	// Преобразуем запрос в сущность
	task := &entity.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	// Если статус не указан, устанавливаем значение по умолчанию
	if task.Status == "" {
		task.Status = entity.StatusTodo
	}

	// Передаем задачу в use case
	if err := h.taskUseCase.CreateTask(r.Context(), task); err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Подготавливаем ответ
	resp := TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetTask обрабатывает запрос на получение задачи по ID
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// В реальном приложении ID будет извлекаться из URL с помощью mux или другого роутера
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	task, err := h.taskUseCase.GetTask(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "access denied") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetAllTasks обрабатывает запрос на получение всех задач пользователя
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskUseCase.GetAllTasks(r.Context())
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var resp []TaskResponse
	for _, task := range tasks {
		resp = append(resp, TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateTask обрабатывает запрос на обновление задачи
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем существующую задачу
	existingTask, err := h.taskUseCase.GetTask(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "access denied") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Обновляем поля
	existingTask.Title = req.Title
	existingTask.Description = req.Description
	if req.Status != "" {
		existingTask.Status = req.Status
	}

	// Сохраняем изменения
	if err := h.taskUseCase.UpdateTask(r.Context(), existingTask); err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "access denied") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := TaskResponse{
		ID:          existingTask.ID,
		Title:       existingTask.Title,
		Description: existingTask.Description,
		Status:      existingTask.Status,
		CreatedAt:   existingTask.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   existingTask.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteTask обрабатывает запрос на удаление задачи
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	if err := h.taskUseCase.DeleteTask(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		if strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "access denied") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
