// internal/router/router.go
package router

import (
	"github.com/SaveljevRoman/go-layout-project-2/internal/handler"
	"net/http"
)

type Router struct {
	Mux         *http.ServeMux // делаем публичным для доступа из main.go
	middlewares []func(http.Handler) http.Handler
}

func NewRouter() *Router {
	return &Router{
		Mux:         http.NewServeMux(),
		middlewares: []func(http.Handler) http.Handler{},
	}
}

// HandleFunc Добавляем метод для регистрации произвольного обработчика
func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.Mux.HandleFunc(pattern, handler)
}

// Остальные методы остаются без изменений

func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) RegisterRoutes(taskHandler *handler.TaskHandler) {
	// Регистрация маршрутов для обработки задач
	r.Mux.HandleFunc("/tasks", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			taskHandler.GetAllTasks(w, req)
		case http.MethodPost:
			taskHandler.CreateTask(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.Mux.HandleFunc("/tasks/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			taskHandler.GetTask(w, req)
		case http.MethodPut:
			taskHandler.UpdateTask(w, req)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (r *Router) Start(addr string) error {
	// Применяем все middleware к mux
	var handler http.Handler = r.Mux
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	// Запускаем сервер
	return http.ListenAndServe(addr, handler)
}
