package db

import "context"

// GatekeeperKind は Gatekeeper チェックの種類。
type GatekeeperKind string

const (
	GateKindMemoryStore  GatekeeperKind = "memory_store"
	GateKindMemoryOutput GatekeeperKind = "memory_output"
)

// GatekeeperProfile はチェック時のポリシープロファイル。
type GatekeeperProfile string

const (
	GateProfileStrict GatekeeperProfile = "STRICT"
	GateProfileNormal GatekeeperProfile = "NORMAL"
	GateProfileDev    GatekeeperProfile = "DEV"
)

// GatekeeperMeta は Gatekeeper 判定に使うメタ情報。
type GatekeeperMeta struct {
	SourceType  string
	SourceTrust string
	Sensitivity string
	Store       StoreKind
}

// GatekeeperCheckRequest は Gatekeeper への入力。
type GatekeeperCheckRequest struct {
	Kind    GatekeeperKind
	Profile GatekeeperProfile
	Content string
	Meta    GatekeeperMeta
}

// GatekeeperDecision は Gatekeeper の判定結果。
type GatekeeperDecision struct {
	Decision   string // "allow" | "deny" | "needs_human"
	Reason     string
	Categories []string
}

// Gatekeeper はメモリ保存／出力時の安全判定を行うインターフェース。
type Gatekeeper interface {
	Check(ctx context.Context, req GatekeeperCheckRequest) (GatekeeperDecision, error)
}
