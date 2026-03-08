package recovery

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRecoveryLogger_PreservesRequiredZeroValueEvidence(t *testing.T) {
	logger := NewRecoveryLogger(t.TempDir())
	if err := logger.Detect("INC-VERIFY-1", []string{"REQ-NFR-002", "REQ-NFR-003", "REQ-NFR-005"}); err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if err := logger.Retry(1); err != nil {
		t.Fatalf("Retry: %v", err)
	}
	if err := logger.Mitigate(0, 1.0); err != nil {
		t.Fatalf("Mitigate: %v", err)
	}
	if err := logger.Resolve(0); err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	summaryPath := filepath.Join(filepath.Dir(filepath.Dir(logger.logPath)), "ops", "incident-summary.json")
	rawSummary, err := os.ReadFile(summaryPath)
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	var summary map[string]any
	if err := json.Unmarshal(rawSummary, &summary); err != nil {
		t.Fatalf("unmarshal summary: %v", err)
	}
	if _, ok := summary["rto_minutes"]; !ok {
		t.Fatalf("expected rto_minutes in summary: %s", rawSummary)
	}
	if _, ok := summary["rpo_minutes"]; !ok {
		t.Fatalf("expected rpo_minutes in summary: %s", rawSummary)
	}
	if got := summary["status"]; got != "resolved" {
		t.Fatalf("expected resolved status, got %#v", got)
	}

	rawLog, err := os.ReadFile(logger.logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	lines := bytes.Split(bytes.TrimSpace(rawLog), []byte("\n"))
	if len(lines) != 4 {
		t.Fatalf("expected 4 log lines, got %d", len(lines))
	}

	var mitigated map[string]any
	if err := json.Unmarshal(lines[2], &mitigated); err != nil {
		t.Fatalf("unmarshal mitigated event: %v", err)
	}
	if _, ok := mitigated["pending_compensation_count"]; !ok {
		t.Fatalf("expected pending_compensation_count in mitigated event: %s", lines[2])
	}
	if got := mitigated["pending_compensation_count"]; got != float64(0) {
		t.Fatalf("expected pending_compensation_count=0, got %#v", got)
	}
	if got := mitigated["short_delete_ready_ratio"]; got != float64(1) {
		t.Fatalf("expected short_delete_ready_ratio=1, got %#v", got)
	}
	if got := mitigated["event"]; got != "mitigated" {
		t.Fatalf("expected mitigated event, got %#v", got)
	}
}

func TestRecoveryLogger_ReplanPersistsTicketID(t *testing.T) {
	logger := NewRecoveryLogger(t.TempDir())
	if err := logger.Detect("INC-VERIFY-2", []string{"REQ-NFR-006"}); err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if err := logger.Replan("IN-2026-03-08-001", "manual follow-up"); err != nil {
		t.Fatalf("Replan: %v", err)
	}

	rawSummary, err := os.ReadFile(logger.sumPath)
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	var summary IncidentSummary
	if err := json.Unmarshal(rawSummary, &summary); err != nil {
		t.Fatalf("unmarshal summary: %v", err)
	}
	if summary.Status != "replan" {
		t.Fatalf("expected replan status, got %q", summary.Status)
	}

	rawLog, err := os.ReadFile(logger.logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	lines := bytes.Split(bytes.TrimSpace(rawLog), []byte("\n"))
	if len(lines) != 2 {
		t.Fatalf("expected 2 log lines, got %d", len(lines))
	}
	var event RecoveryEvent
	if err := json.Unmarshal(lines[1], &event); err != nil {
		t.Fatalf("unmarshal event: %v", err)
	}
	if event.Event != "replan" {
		t.Fatalf("expected replan event, got %q", event.Event)
	}
	if event.ReplanTicketID != "IN-2026-03-08-001" {
		t.Fatalf("expected replan ticket id, got %q", event.ReplanTicketID)
	}
}
