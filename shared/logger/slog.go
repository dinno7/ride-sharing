package logger

import (
	"log/slog"
	"os"
)

type slogLogger struct {
	l *slog.Logger
}

func NewSlogLogger(l *slog.Logger) Logger {
	return &slogLogger{l: l}
}

func (s *slogLogger) Debug(msg string, args ...any) {
	s.l.Debug(msg, args...)
}

func (s *slogLogger) Info(msg string, args ...any) {
	s.l.Info(msg, args...)
}

func (s *slogLogger) Warn(msg string, args ...any) {
	s.l.Warn(msg, args...)
}

func (s *slogLogger) Error(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, slog.Any("error", err))
	}
	s.l.Error(msg, args...)
}

func (s *slogLogger) Fatal(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, slog.Any("error", err))
	}
	s.l.Error(msg, args...)
	os.Exit(1)
}

func (s *slogLogger) With(args ...any) Logger {
	return &slogLogger{l: s.l.With(args...)}
}
