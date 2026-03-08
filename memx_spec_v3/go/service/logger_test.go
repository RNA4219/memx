package service

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
)

var errTestLLM = errors.New("llm unavailable")

func newBufferedTestLogger() (*bytes.Buffer, warningLogger) {
	buf := &bytes.Buffer{}
	return buf, log.New(buf, "", 0)
}

func assertAutoSummaryWarningLogged(t *testing.T, output, store, title string, err error) {
	t.Helper()

	checks := []string{
		"WARN auto-summary failed",
		"store=" + store,
		"title=" + `"` + title + `"`,
		"err=" + err.Error(),
	}
	for _, want := range checks {
		if !strings.Contains(output, want) {
			t.Fatalf("expected warning log to contain %q, got %q", want, output)
		}
	}
}
