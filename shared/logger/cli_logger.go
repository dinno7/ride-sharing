package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type cliLogger struct {
	prefix string
	logger *log.Logger
}

func NewCliLogger(prefix string) Logger {
	prefix = fmt.Sprintf("%s ", strings.TrimSpace(prefix))
	return &cliLogger{
		logger: log.New(os.Stdout, prefix, log.LstdFlags),
	}
}

func (l *cliLogger) With(args ...any) Logger {
	newPrefix := l.prefix
	if len(args) > 0 {
		parts := make([]string, 0, len(args))
		for i := 0; i < len(args); i++ {
			parts = append(parts, fmt.Sprint(args[i]))
		}
		ctx := strings.Join(parts, " ")
		if newPrefix != "" {
			newPrefix += " "
		}
		newPrefix += ctx
	}
	return &cliLogger{
		prefix: newPrefix,
		logger: l.logger,
	}
}

// ANSI Color Constants
type colorCode string

const (
	colorReset  colorCode = "\033[0m"
	colorRed    colorCode = "\033[31m"
	colorGreen  colorCode = "\033[32m"
	colorYellow colorCode = "\033[33m"
	colorCyan   colorCode = "\033[36m"
	colorPurple colorCode = "\033[35m"

	// Bold Colors (prefix with 1;)
	colorBoldRed    = "\033[1;31m"
	colorBoldPurple = "\033[1;35m"
)

func (l *cliLogger) Debug(msg string, args ...any) {
	l.printf(colorText(colorCyan, "[DEBUG]"), colorText(colorCyan, msg), args...)
}

func (l *cliLogger) Info(msg string, args ...any) {
	l.printf(colorText(colorGreen, "[INFO]"), colorText(colorGreen, msg), args...)
}

func (l *cliLogger) Warn(msg string, args ...any) {
	l.printf(colorText(colorYellow, "[WARN]"), colorText(colorYellow, msg), args...)
}

func (l *cliLogger) Error(msg string, err error, args ...any) {
	fullMsg := colorText(colorBoldRed, fmt.Sprintf("%s (error: %v)", msg, err))
	l.printf(colorText(colorRed, "[ERROR]"), fullMsg, args...)
}

func (l *cliLogger) Fatal(msg string, err error, args ...any) {
	fullMsg := colorText(colorBoldPurple, fmt.Sprintf("%s (error: %v)", msg, err))
	l.printf(colorText(colorPurple, "[FATAL]"), fullMsg, args...)
	os.Exit(1)
}

// Helper for formatted output
func (l *cliLogger) printf(level string, msg string, args ...any) {
	text := fmt.Sprintf(msg, args...)
	if l.prefix != "" {
		text = fmt.Sprintf("%s | %s", l.prefix, text)
	}
	l.logger.Printf("%s %s", level, text)
}

func colorText(color colorCode, msg string) string {
	return fmt.Sprintf("%s%s%s", color, msg, colorReset)
}
