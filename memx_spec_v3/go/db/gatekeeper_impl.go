package db

import (
	"context"
	"strings"
)

// Decision 定数
const (
	DecisionAllow      = "allow"
	DecisionDeny       = "deny"
	DecisionNeedsHuman = "needs_human"
)

// DefaultGatekeeper はデフォルトの Gatekeeper 実装。
// ポリシープロファイルに基づいてシンプルなルールベース判定を行う。
type DefaultGatekeeper struct {
	profile GatekeeperProfile
}

// NewDefaultGatekeeper は DefaultGatekeeper を作成する。
func NewDefaultGatekeeper(profile GatekeeperProfile) *DefaultGatekeeper {
	if profile == "" {
		profile = GateProfileNormal
	}
	return &DefaultGatekeeper{profile: profile}
}

// Check は Gatekeeper チェックを実行する。
// ルール:
// - sensitivity=secret の場合は常に deny（fail-closed）
// - DEV プロファイルでは secret 以外は allow
// - NORMAL プロファイルでは source_trust=untrusted の場合は needs_human
// - STRICT プロファイルでは source_trust が trusted でない限り needs_human
func (g *DefaultGatekeeper) Check(ctx context.Context, req GatekeeperCheckRequest) (GatekeeperDecision, error) {
	// コンテキストのキャンセルをチェック
	select {
	case <-ctx.Done():
		return GatekeeperDecision{}, ctx.Err()
	default:
	}

	// sensitivity の正規化
	sensitivity := strings.ToLower(strings.TrimSpace(req.Meta.Sensitivity))

	// secret は常に deny（fail-closed: 安全側に倒す）
	if sensitivity == "secret" {
		return GatekeeperDecision{
			Decision:   DecisionDeny,
			Reason:     "content marked as secret cannot be stored",
			Categories: []string{"sensitivity:secret"},
		}, nil
	}

	// プロファイル別の判定
	switch g.profile {
	case GateProfileDev:
		// DEV: secret 以外は allow（開発環境向け）
		return GatekeeperDecision{
			Decision:   DecisionAllow,
			Reason:     "DEV profile allows non-secret content",
			Categories: nil,
		}, nil

	case GateProfileStrict:
		// STRICT: source_trust が trusted でない限り needs_human
		sourceTrust := strings.ToLower(strings.TrimSpace(req.Meta.SourceTrust))
		if sourceTrust != "trusted" {
			return GatekeeperDecision{
				Decision:   DecisionNeedsHuman,
				Reason:     "STRICT profile requires human review for non-trusted sources",
				Categories: []string{"source_trust:" + sourceTrust},
			}, nil
		}
		return GatekeeperDecision{
			Decision:   DecisionAllow,
			Reason:     "trusted source approved under STRICT profile",
			Categories: nil,
		}, nil

	default: // GateProfileNormal
		// NORMAL: source_trust=untrusted の場合は needs_human
		sourceTrust := strings.ToLower(strings.TrimSpace(req.Meta.SourceTrust))
		if sourceTrust == "untrusted" {
			return GatekeeperDecision{
				Decision:   DecisionNeedsHuman,
				Reason:     "untrusted source requires human review under NORMAL profile",
				Categories: []string{"source_trust:untrusted"},
			}, nil
		}
		return GatekeeperDecision{
			Decision:   DecisionAllow,
			Reason:     "NORMAL profile allows content",
			Categories: nil,
		}, nil
	}
}

// AllowAllGatekeeper は全てのリクエストを許可する Gatekeeper（テスト用）。
type AllowAllGatekeeper struct{}

func (g *AllowAllGatekeeper) Check(ctx context.Context, req GatekeeperCheckRequest) (GatekeeperDecision, error) {
	return GatekeeperDecision{
		Decision:   DecisionAllow,
		Reason:     "allow-all gatekeeper",
		Categories: nil,
	}, nil
}

// DenyAllGatekeeper は全てのリクエストを拒否する Gatekeeper（テスト用）。
type DenyAllGatekeeper struct{}

func (g *DenyAllGatekeeper) Check(ctx context.Context, req GatekeeperCheckRequest) (GatekeeperDecision, error) {
	return GatekeeperDecision{
		Decision:   DecisionDeny,
		Reason:     "deny-all gatekeeper",
		Categories: []string{"deny-all"},
	}, nil
}