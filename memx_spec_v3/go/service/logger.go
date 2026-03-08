package service

import (
	"log"
	"os"
	"strings"
)

type warningLogger interface {
	Printf(format string, v ...any)
}

func newDefaultLogger() warningLogger {
	return log.New(os.Stderr, "memx: ", 0)
}

// SetLogger は警告ログ出力先を差し替える。
func (s *Service) SetLogger(logger warningLogger) {
	s.Logger = logger
}

func (s *Service) warnf(format string, args ...any) {
	if s == nil || s.Logger == nil {
		return
	}
	s.Logger.Printf("WARN "+format, args...)
}

func (s *Service) warnAutoSummaryFailure(store, title string, err error) {
	if err == nil {
		return
	}
	s.warnf("auto-summary failed store=%s title=%q err=%v", store, compactLogValue(title, 80), err)
}

func compactLogValue(v string, limit int) string {
	trimmed := strings.TrimSpace(v)
	if limit <= 0 || len(trimmed) <= limit {
		return trimmed
	}
	if limit <= 3 {
		return trimmed[:limit]
	}
	return trimmed[:limit-3] + "..."
}
