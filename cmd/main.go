// cmd/main.go (альтернативная версия)
package main

import (
	"github.com/SaveljevRoman/go-layout-project-2/internal/adapter"
	"github.com/SaveljevRoman/go-layout-project-2/internal/handler"
	"github.com/SaveljevRoman/go-layout-project-2/internal/repository/db"
	"github.com/SaveljevRoman/go-layout-project-2/internal/router"
	"github.com/SaveljevRoman/go-layout-project-2/internal/usecase"
	dbpkg "github.com/SaveljevRoman/go-layout-project-2/pkg/db"
	"github.com/SaveljevRoman/go-layout-project-2/pkg/logger"
	"log"
	"net/http"
)

func main() {
	// Инициализация логгера
	appLogger := logger.NewLogger()

	// Конфигурация базы данных
	dbConfig := dbpkg.DBConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "user",
		Password: "password",
		DBName:   "taskmanager",
	}

	// Инициализация базы данных
	database := dbpkg.NewDatabase(dbConfig)
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Инициализация хранилищ
	taskRepo := db.NewTaskRepository()

	// Инициализация use cases
	taskUseCase := usecase.NewTaskUseCase(taskRepo, appLogger)

	// Инициализация адаптеров
	taskAPI := adapter.NewTaskAPI(taskUseCase)

	// Инициализация обработчиков
	taskHandler := handler.NewTaskHandler(taskUseCase)

	// Инициализация роутера
	r := router.NewRouter()

	// Регистрация middleware
	r.Use(router.LoggingMiddleware(appLogger))
	r.Use(router.AuthMiddleware())

	// Регистрация маршрутов
	r.RegisterRoutes(taskHandler)

	// Используем новый метод HandleFunc для регистрации обработчика API
	r.HandleFunc("/api/export", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data, err := taskAPI.ExportTasksToJSON(req.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=tasks.json")
		w.Write(data)
	})

	// Запуск сервера
	log.Println("Starting server on :8080")
	if err := r.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
