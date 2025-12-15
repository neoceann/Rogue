package logger

import (
	"fmt"
	"strings"
	"time"
)

type BufferLogger struct {
	buf     strings.Builder
	start   time.Time
	verbose bool // debug on/off
}

func NewBufferLogger(start time.Time, verbose bool) *BufferLogger {
	return &BufferLogger{
		start:   start,
		verbose: verbose,
	}
}

func (l *BufferLogger) log(level, format string, args ...any) {
	// if !l.verbose && level == "DEBUG" {
	if !l.verbose {
		return
	}
	elapsed := time.Since(l.start).Seconds()
	msg := fmt.Sprintf("[%.3f][%s] %s\n",
		elapsed, level, fmt.Sprintf(format, args...))
	l.buf.WriteString(msg)
}

func (l *BufferLogger) Debug(format string, args ...any) { l.log("DEBUG", format, args...) }
func (l *BufferLogger) Info(format string, args ...any)  { l.log("INFO ", format, args...) }

// Bytes возвращает лог как []byte — удобно для записи в файл или stderr
func (l *BufferLogger) Bytes() []byte {
	return []byte(l.buf.String())
}

// Reset очищает буфер (если захочешь логировать по уровням отдельно)
func (l *BufferLogger) Reset() {
	l.buf.Reset()
}
