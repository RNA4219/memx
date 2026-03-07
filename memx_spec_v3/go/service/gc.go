package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GCConfig は GC の設定値。
// requirements-cli.md#3-5 の閾値と一致させる。
type GCConfig struct {
	SoftLimitNotes      int // 1200
	HardLimitNotes      int // 2000
	MinIntervalMinutes  int // 180
	ArchiveBatchSize    int // 100 (1回のGCで退避する最大ノート数)
}

// DefaultGCConfig はデフォルトのGC設定。
var DefaultGCConfig = GCConfig{
	SoftLimitNotes:     1200,
	HardLimitNotes:     2000,
	MinIntervalMinutes: 180,
	ArchiveBatchSize:   100,
}

// GCTriggerDecision は Phase0 の判定結果。
type GCTriggerDecision struct {
	ShouldRun bool   `json:"should_run"`
	Reason    string `json:"reason"`
	Metrics   struct {
		NoteCount           int `json:"note_count"`
		SoftLimitNotes      int `json:"soft_limit_notes"`
		HardLimitNotes      int `json:"hard_limit_notes"`
		MinutesSinceLastGC  int `json:"minutes_since_last_gc"`
		MinIntervalMinutes  int `json:"min_interval_minutes"`
	} `json:"metrics"`
}

// GCDryRunResult は dry-run の結果。
type GCDryRunResult struct {
	Target     string             `json:"target"`
	Phase      string             `json:"phase"`
	Decision   GCTriggerDecision  `json:"decision"`
	PlannedOps []GCPlannedOp      `json:"planned_ops"`
}

// GCPlannedOp は予定されている操作。
type GCPlannedOp struct {
	Op               string   `json:"op"`
	SrcNoteIDs       []string `json:"src_note_ids,omitempty"`
	SrcNoteID        string   `json:"src_note_id,omitempty"`
	DestStore        string   `json:"dest_store"`
	LineageRelation  string   `json:"lineage_relation,omitempty"`
}

// GCRequest は GC 実行リクエスト。
type GCRequest struct {
	Target  string
	DryRun  bool
	Enabled bool // feature flag
}

// GCResult は GC 実行結果。
type GCResult struct {
	Status     string         `json:"status"`
	DryRun     bool           `json:"dry_run"`
	DryRunResult *GCDryRunResult `json:"dry_run_result,omitempty"`
	ArchivedCount int          `json:"archived_count,omitempty"`
}

// GCShort は short.db に対する GC を実行する。
// v1.x では Phase 0 (トリガ判定) と Phase 3 (Archive退避) のみ実装。
func (s *Service) GCShort(ctx context.Context, req GCRequest) (GCResult, error) {
	if !req.Enabled {
		return GCResult{}, ErrFeatureDisabled
	}

	cfg := DefaultGCConfig

	// Phase 0: トリガ判定
	decision, err := s.gcTriggerCheck(ctx, cfg)
	if err != nil {
		return GCResult{}, fmt.Errorf("trigger check: %w", err)
	}

	// dry-run の場合は判定結果のみ返す
	if req.DryRun {
		result := s.gcBuildDryRunResult(decision, cfg)
		return GCResult{
			Status:       "dry_run",
			DryRun:       true,
			DryRunResult: result,
		}, nil
	}

	// 実行不要な場合はスキップ
	if !decision.ShouldRun {
		return GCResult{
			Status: "skipped",
			DryRun: false,
		}, nil
	}

	// Phase 3: Archive退避（簡易実装）
	archivedCount, err := s.gcArchiveOldNotes(ctx, cfg)
	if err != nil {
		return GCResult{}, fmt.Errorf("archive: %w", err)
	}

	// last_gc_at 更新
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, _ = s.Conn.DB.ExecContext(ctx,
		`UPDATE short_meta SET value = ? WHERE key = 'last_gc_at'`, now)

	return GCResult{
		Status:        "ok",
		DryRun:        false,
		ArchivedCount: archivedCount,
	}, nil
}

// gcTriggerCheck は GC 実行要否を判定する。
func (s *Service) gcTriggerCheck(ctx context.Context, cfg GCConfig) (GCTriggerDecision, error) {
	var decision GCTriggerDecision
	decision.Metrics.SoftLimitNotes = cfg.SoftLimitNotes
	decision.Metrics.HardLimitNotes = cfg.HardLimitNotes
	decision.Metrics.MinIntervalMinutes = cfg.MinIntervalMinutes

	// ノート数取得
	var noteCount int
	err := s.Conn.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notes`).Scan(&noteCount)
	if err != nil {
		return decision, fmt.Errorf("count notes: %w", err)
	}
	decision.Metrics.NoteCount = noteCount

	// last_gc_at 取得
	var lastGCAt string
	err = s.Conn.DB.QueryRowContext(ctx,
		`SELECT value FROM short_meta WHERE key = 'last_gc_at'`).Scan(&lastGCAt)
	if err != nil && err != sql.ErrNoRows {
		return decision, fmt.Errorf("get last_gc_at: %w", err)
	}

	// 経過時間計算
	minutesSinceLastGC := 0
	if lastGCAt != "" && lastGCAt != "1970-01-01T00:00:00Z" {
		lastGC, err := time.Parse(time.RFC3339Nano, lastGCAt)
		if err == nil {
			minutesSinceLastGC = int(time.Since(lastGC).Minutes())
		}
	}
	decision.Metrics.MinutesSinceLastGC = minutesSinceLastGC

	// 判定ロジック
	if noteCount >= cfg.HardLimitNotes {
		decision.ShouldRun = true
		decision.Reason = "hard_limit_reached"
	} else if noteCount >= cfg.SoftLimitNotes && minutesSinceLastGC >= cfg.MinIntervalMinutes {
		decision.ShouldRun = true
		decision.Reason = "soft_limit_reached"
	} else if noteCount >= cfg.SoftLimitNotes {
		decision.ShouldRun = false
		decision.Reason = "interval_not_elapsed"
	} else {
		decision.ShouldRun = false
		decision.Reason = "under_limit"
	}

	return decision, nil
}

// gcBuildDryRunResult は dry-run 結果を構築する。
func (s *Service) gcBuildDryRunResult(decision GCTriggerDecision, cfg GCConfig) *GCDryRunResult {
	result := &GCDryRunResult{
		Target:   "short",
		Phase:    "phase0",
		Decision: decision,
	}

	// 実行予定がある場合は操作を追加
	if decision.ShouldRun {
		// 簡易実装: 古いノートを退避予定として表示
		result.PlannedOps = append(result.PlannedOps, GCPlannedOp{
			Op:              "archive_candidates",
			DestStore:       "archive",
			LineageRelation: "archived_from",
		})
	}

	return result
}

// gcArchiveOldNotes は古いノートを archive へ退避する。
// v1.x 簡易実装: アクセス数が少なく古いノートを退避。
func (s *Service) gcArchiveOldNotes(ctx context.Context, cfg GCConfig) (int, error) {
	// 退避候補を取得（アクセス数0、かつ作成から30日以上経過）
	cutoffDate := time.Now().AddDate(0, 0, -30).Format(time.RFC3339Nano)

	rows, err := s.Conn.DB.QueryContext(ctx, `
SELECT id, title, summary, body, created_at, updated_at, last_accessed_at,
       access_count, source_type, origin, source_trust, sensitivity
FROM notes
WHERE access_count = 0 AND created_at < ?
ORDER BY created_at ASC
LIMIT ?`, cutoffDate, cfg.ArchiveBatchSize)
	if err != nil {
		return 0, fmt.Errorf("query archive candidates: %w", err)
	}
	defer rows.Close()

	var candidates []struct {
		id, title, summary, body, createdAt, updatedAt, lastAccessedAt string
		accessCount                                                     int
		sourceType, origin, sourceTrust, sensitivity                    string
	}

	for rows.Next() {
		var c struct {
			id, title, summary, body, createdAt, updatedAt, lastAccessedAt string
			accessCount                                                     int
			sourceType, origin, sourceTrust, sensitivity                    string
		}
		if err := rows.Scan(&c.id, &c.title, &c.summary, &c.body,
			&c.createdAt, &c.updatedAt, &c.lastAccessedAt,
			&c.accessCount, &c.sourceType, &c.origin, &c.sourceTrust, &c.sensitivity); err != nil {
			return 0, fmt.Errorf("scan candidate: %w", err)
		}
		candidates = append(candidates, c)
	}

	if len(candidates) == 0 {
		return 0, nil
	}

	// Archive へ退避
	tx, err := s.Conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now().UTC().Format(time.RFC3339Nano)
	archivedCount := 0

	for _, c := range candidates {
		// archive.notes へ INSERT
		_, err := tx.ExecContext(ctx, `
INSERT OR IGNORE INTO notes(
  id, title, summary, body, created_at, updated_at, last_accessed_at,
  access_count, source_type, origin, source_trust, sensitivity
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			c.id, c.title, c.summary, c.body,
			c.createdAt, c.updatedAt, c.lastAccessedAt,
			c.accessCount, c.sourceType, c.origin, c.sourceTrust, c.sensitivity)
		if err != nil {
			continue // エラーはスキップ
		}

		// lineage へ記録
		_, err = tx.ExecContext(ctx, `
INSERT INTO lineage(src_store, src_note_id, dest_store, dest_note_id, relation, created_at)
VALUES (?, ?, ?, ?, ?, ?)`,
			"short", c.id, "archive", c.id, "archived_from", now)
		if err != nil {
			continue
		}

		archivedCount++
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}

	return archivedCount, nil
}

// ToJSON は dry-run 結果を JSON 文字列で返す。
func (r *GCDryRunResult) ToJSON() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

// FormatDryRunOutput は dry-run 結果を人間可読形式で返す。
func (r *GCDryRunResult) FormatDryRunOutput() string {
	var sb strings.Builder
	sb.WriteString("GC Dry-Run Result:\n")
	sb.WriteString(fmt.Sprintf("  Target: %s\n", r.Target))
	sb.WriteString(fmt.Sprintf("  Phase: %s\n", r.Phase))
	sb.WriteString(fmt.Sprintf("  Should Run: %v\n", r.Decision.ShouldRun))
	sb.WriteString(fmt.Sprintf("  Reason: %s\n", r.Decision.Reason))
	sb.WriteString("  Metrics:\n")
	sb.WriteString(fmt.Sprintf("    Note Count: %d\n", r.Decision.Metrics.NoteCount))
	sb.WriteString(fmt.Sprintf("    Soft Limit: %d\n", r.Decision.Metrics.SoftLimitNotes))
	sb.WriteString(fmt.Sprintf("    Hard Limit: %d\n", r.Decision.Metrics.HardLimitNotes))
	sb.WriteString(fmt.Sprintf("    Minutes Since Last GC: %d\n", r.Decision.Metrics.MinutesSinceLastGC))
	sb.WriteString(fmt.Sprintf("    Min Interval (minutes): %d\n", r.Decision.Metrics.MinIntervalMinutes))

	if len(r.PlannedOps) > 0 {
		sb.WriteString("  Planned Operations:\n")
		for _, op := range r.PlannedOps {
			sb.WriteString(fmt.Sprintf("    - %s -> %s\n", op.Op, op.DestStore))
		}
	}

	return sb.String()
}