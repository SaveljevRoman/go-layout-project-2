package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors содержит список ошибок валидации
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// ValidateRequest Функция для валидации запроса
func ValidateRequest(r *http.Request, schema interface{}) ([]ValidationError, error) {
	var validationErrors []ValidationError

	if err := json.NewDecoder(r.Body).Decode(schema); err != nil {
		return nil, err
	}

	// Здесь мы можем использовать различные валидационные библиотеки
	// В этом примере используем простую проверку структуры

	switch s := schema.(type) {
	case *CreateTaskRequest:
		if s.Title == "" {
			validationErrors = append(validationErrors, ValidationError{
				Field:   "title",
				Message: "Title is required",
			})
		}

		if len(s.Title) > 100 {
			validationErrors = append(validationErrors, ValidationError{
				Field:   "title",
				Message: "Title must be less than 100 characters",
			})
		}

		if len(s.Description) > 1000 {
			validationErrors = append(validationErrors, ValidationError{
				Field:   "description",
				Message: "Description must be less than 1000 characters",
			})
		}

		if s.Status != "" {
			status := strings.ToUpper(string(s.Status))
			if status != "TODO" && status != "IN_PROGRESS" && status != "DONE" {
				validationErrors = append(validationErrors, ValidationError{
					Field:   "status",
					Message: "Status must be one of: TODO, IN_PROGRESS, DONE",
				})
			}
		}
	}

	return validationErrors, nil
}
