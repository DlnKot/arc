package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

type FileLogger struct {
	mu     sync.Mutex
	logger *log.Logger
	file   *os.File
}

func New(appName string) (*FileLogger, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	logDir := filepath.Join(configDir, appName, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	return &FileLogger{
		logger: log.New(file, "", 0),
		file:   file,
	}, nil
}

func (l *FileLogger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	return l.file.Close()
}

func (l *FileLogger) Infof(format string, args ...any) {
	l.log("info", format, args...)
}

func (l *FileLogger) Warnf(format string, args ...any) {
	l.log("warn", format, args...)
}

func (l *FileLogger) Errorf(format string, args ...any) {
	l.log("error", format, args...)
}

func (l *FileLogger) log(level string, format string, args ...any) {
	if l == nil || l.logger == nil {
		return
	}

	message := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("[%s] [%s] %s", time.Now().Format(time.RFC3339), strings.ToUpper(level), message)

	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Println(line)
	fmt.Println(line)
}
