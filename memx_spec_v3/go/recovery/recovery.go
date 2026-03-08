package recovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// IncidentSummary はインシデント概要の正本証跡。
// 要件: REQ-NFR-002, REQ-NFR-003, REQ-NFR-005
type IncidentSummary struct {
	IncidentID          string   `json:"incident_id"`
	DetectedAt          string   `json:"detected_at"`
	MitigatedAt         string   `json:"mitigated_at,omitempty"`
	ResolvedAt          string   `json:"resolved_at,omitempty"`
	RTOMinutes          *int     `json:"rto_minutes,omitempty"`
	RPOMinutes          *int     `json:"rpo_minutes,omitempty"`
	Status              string   `json:"status"` // detected, mitigated, resolved, replan
	RelatedRequirements []string `json:"related_requirements,omitempty"`
	Waiver              *Waiver  `json:"waiver,omitempty"`
	RecoveryLogPath     string   `json:"recovery_log_path,omitempty"`
	IncidentRecordPath  string   `json:"incident_record_path,omitempty"`
}

// Waiver は要件免除情報。
type Waiver struct {
	RequirementID      string `json:"requirement_id"`
	Reason             string `json:"reason"`
	ExpiresAt          string `json:"expires_at,omitempty"`
	Approver           string `json:"approver,omitempty"`
	AlternativeControl string `json:"alternative_control,omitempty"`
	LiftConditions     string `json:"lift_conditions,omitempty"`
}

// RecoveryEvent は復旧ログの1イベント。
// NDJSON形式で追記される。
type RecoveryEvent struct {
	Timestamp                string   `json:"timestamp"`
	Event                    string   `json:"event"` // detect, retry, rollback, mitigated, resolved, replan
	IncidentID               string   `json:"incident_id,omitempty"`
	RetryCount               *int     `json:"retry_count,omitempty"`
	PendingCompensationCount *int     `json:"pending_compensation_count,omitempty"`
	ShortDeleteReadyRatio    *float64 `json:"short_delete_ready_ratio,omitempty"`
	ReplanTicketID           string   `json:"replan_ticket_id,omitempty"`
	Message                  string   `json:"message,omitempty"`
}

// RecoveryLogger は復旧ログを管理する。
type RecoveryLogger struct {
	mu       sync.Mutex
	logPath  string
	sumPath  string
	incident *IncidentSummary
}

// NewRecoveryLogger は新しい復旧ロガーを作成する。
func NewRecoveryLogger(artifactsDir string) *RecoveryLogger {
	return &RecoveryLogger{
		logPath: filepath.Join(artifactsDir, "ops", "recovery-log.ndjson"),
		sumPath: filepath.Join(artifactsDir, "ops", "incident-summary.json"),
	}
}

// Detect は障害検知を記録する。
func (l *RecoveryLogger) Detect(incidentID string, requirements []string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)

	l.incident = &IncidentSummary{
		IncidentID:          incidentID,
		DetectedAt:          now,
		Status:              "detected",
		RelatedRequirements: requirements,
		RecoveryLogPath:     l.logPath,
	}

	if err := l.saveSummary(); err != nil {
		return err
	}

	return l.appendEvent(&RecoveryEvent{
		Timestamp:  now,
		Event:      "detect",
		IncidentID: incidentID,
		Message:    "Incident detected",
	})
}

// Retry は再試行を記録する。
func (l *RecoveryLogger) Retry(retryCount int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	retryCountCopy := retryCount
	return l.appendEvent(&RecoveryEvent{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Event:      "retry",
		RetryCount: &retryCountCopy,
	})
}

// Mitigate は暫定復旧を記録する。
func (l *RecoveryLogger) Mitigate(pendingCompensation int, deleteReadyRatio float64) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.incident == nil {
		return nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	l.incident.MitigatedAt = now
	l.incident.Status = "mitigated"

	// RTO計算
	detected, err := time.Parse(time.RFC3339, l.incident.DetectedAt)
	if err == nil {
		rtoMinutes := int(time.Since(detected).Minutes())
		l.incident.RTOMinutes = &rtoMinutes
	}

	if err := l.saveSummary(); err != nil {
		return err
	}

	pendingCompensationCopy := pendingCompensation
	deleteReadyRatioCopy := deleteReadyRatio
	return l.appendEvent(&RecoveryEvent{
		Timestamp:                now,
		Event:                    "mitigated",
		PendingCompensationCount: &pendingCompensationCopy,
		ShortDeleteReadyRatio:    &deleteReadyRatioCopy,
	})
}

// Resolve は恒久復旧を記録する。
func (l *RecoveryLogger) Resolve(rpoMinutes int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.incident == nil {
		return nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	l.incident.ResolvedAt = now
	l.incident.Status = "resolved"
	rpoMinutesCopy := rpoMinutes
	l.incident.RPOMinutes = &rpoMinutesCopy

	// RTO再計算（恒久復旧ベース）
	detected, err := time.Parse(time.RFC3339, l.incident.DetectedAt)
	if err == nil {
		rtoMinutes := int(time.Since(detected).Minutes())
		l.incident.RTOMinutes = &rtoMinutes
	}

	if err := l.saveSummary(); err != nil {
		return err
	}

	return l.appendEvent(&RecoveryEvent{
		Timestamp: now,
		Event:     "resolved",
	})
}

// Replan は再計画への移行を記録する。
func (l *RecoveryLogger) Replan(ticketID, reason string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.incident == nil {
		return nil
	}

	l.incident.Status = "replan"

	if err := l.saveSummary(); err != nil {
		return err
	}

	return l.appendEvent(&RecoveryEvent{
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		Event:          "replan",
		ReplanTicketID: ticketID,
		Message:        reason,
	})
}

// Rollback はロールバックを記録する。
func (l *RecoveryLogger) Rollback(reason string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.appendEvent(&RecoveryEvent{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Event:     "rollback",
		Message:   reason,
	})
}

// GetSummary は現在のインシデント概要を返す。
func (l *RecoveryLogger) GetSummary() *IncidentSummary {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.incident
}

func (l *RecoveryLogger) saveSummary() error {
	dir := filepath.Dir(l.sumPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(l.incident, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(l.sumPath, data, 0644)
}

func (l *RecoveryLogger) appendEvent(event *RecoveryEvent) error {
	dir := filepath.Dir(l.logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(l.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = f.Write(append(data, '\n'))
	return err
}
