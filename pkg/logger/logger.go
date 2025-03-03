package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log(l.infoLogger, message, fields)
}

func (l *Logger) Error(message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	l.log(l.errorLogger, message, fields)
}

func (l *Logger) log(logger *log.Logger, message string, fields map[string]interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	fieldsStr := ""

	for k, v := range fields {
		fieldsStr += fmt.Sprintf("%s=%v ", k, v)
	}

	logger.Printf("[%s] %s %s\n", timestamp, message, fieldsStr)
}
