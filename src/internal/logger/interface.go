package logger

// Logger — интерфейс для логгирования.
type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
}

// NoopLogger — "тихий" логгер (для тестов или release без логов).
type NoopLogger struct{}

func (NoopLogger) Debug(string, ...any) {}
func (NoopLogger) Info(string, ...any)  {}
