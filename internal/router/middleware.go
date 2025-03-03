package router

import (
	"context"
	"github.com/SaveljevRoman/go-layout-project-2/pkg/logger"
	"net/http"
	"strings"
	"time"
)

// LoggingMiddleware логирует информацию о запросе
func LoggingMiddleware(logger *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Логируем входящий запрос
			logger.Info("Request started", map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
				"remote": r.RemoteAddr,
			})

			// Вызываем следующий обработчик
			next.ServeHTTP(w, r)

			// Логируем время выполнения запроса
			logger.Info("Request completed", map[string]interface{}{
				"method":   r.Method,
				"path":     r.URL.Path,
				"duration": time.Since(start).String(),
			})
		})
	}
}

// AuthMiddleware проверяет авторизацию пользователя
func AuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Проверяем формат токена
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			// Извлекаем токен
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// В реальном приложении здесь была бы проверка токена
			// Для примера просто проверяем, что токен не пустой
			if token == "" {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Предположим, что из токена мы получили userID
			userID := "user-123" // В реальности это значение получается из JWT токена

			// Добавляем userID в контекст запроса
			ctx := context.WithValue(r.Context(), "user_id", userID)

			// Вызываем следующий обработчик с обновленным контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
