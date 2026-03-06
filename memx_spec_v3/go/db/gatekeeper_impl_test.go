package db

import (
	"context"
	"testing"
)

func TestDefaultGatekeeper_SecretDeny(t *testing.T) {
	gk := NewDefaultGatekeeper(GateProfileNormal)

	req := GatekeeperCheckRequest{
		Kind:    GateKindMemoryStore,
		Profile: GateProfileNormal,
		Content: "test content",
		Meta: GatekeeperMeta{
			Sensitivity: "secret",
		},
	}

	decision, err := gk.Check(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if decision.Decision != DecisionDeny {
		t.Errorf("expected deny for secret, got: %s", decision.Decision)
	}
}

func TestDefaultGatekeeper_DevProfile(t *testing.T) {
	gk := NewDefaultGatekeeper(GateProfileDev)

	tests := []struct {
		name        string
		sensitivity string
		sourceTrust string
		want        string
	}{
		{"secret always deny", "secret", "trusted", DecisionDeny},
		{"public allow", "public", "untrusted", DecisionAllow},
		{"internal allow", "internal", "user_input", DecisionAllow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := GatekeeperCheckRequest{
				Kind:    GateKindMemoryStore,
				Profile: GateProfileDev,
				Meta: GatekeeperMeta{
					Sensitivity: tt.sensitivity,
					SourceTrust: tt.sourceTrust,
				},
			}

			decision, err := gk.Check(context.Background(), req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if decision.Decision != tt.want {
				t.Errorf("expected %s, got: %s", tt.want, decision.Decision)
			}
		})
	}
}

func TestDefaultGatekeeper_NormalProfile(t *testing.T) {
	gk := NewDefaultGatekeeper(GateProfileNormal)

	tests := []struct {
		name        string
		sensitivity string
		sourceTrust string
		want        string
	}{
		{"secret deny", "secret", "trusted", DecisionDeny},
		{"untrusted needs human", "public", "untrusted", DecisionNeedsHuman},
		{"user_input allow", "internal", "user_input", DecisionAllow},
		{"trusted allow", "public", "trusted", DecisionAllow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := GatekeeperCheckRequest{
				Kind:    GateKindMemoryStore,
				Profile: GateProfileNormal,
				Meta: GatekeeperMeta{
					Sensitivity: tt.sensitivity,
					SourceTrust: tt.sourceTrust,
				},
			}

			decision, err := gk.Check(context.Background(), req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if decision.Decision != tt.want {
				t.Errorf("expected %s, got: %s", tt.want, decision.Decision)
			}
		})
	}
}

func TestDefaultGatekeeper_StrictProfile(t *testing.T) {
	gk := NewDefaultGatekeeper(GateProfileStrict)

	tests := []struct {
		name        string
		sensitivity string
		sourceTrust string
		want        string
	}{
		{"secret deny", "secret", "trusted", DecisionDeny},
		{"untrusted needs human", "public", "untrusted", DecisionNeedsHuman},
		{"user_input needs human", "internal", "user_input", DecisionNeedsHuman},
		{"trusted allow", "public", "trusted", DecisionAllow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := GatekeeperCheckRequest{
				Kind:    GateKindMemoryStore,
				Profile: GateProfileStrict,
				Meta: GatekeeperMeta{
					Sensitivity: tt.sensitivity,
					SourceTrust: tt.sourceTrust,
				},
			}

			decision, err := gk.Check(context.Background(), req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if decision.Decision != tt.want {
				t.Errorf("expected %s, got: %s", tt.want, decision.Decision)
			}
		})
	}
}

func TestAllowAllGatekeeper(t *testing.T) {
	gk := &AllowAllGatekeeper{}

	req := GatekeeperCheckRequest{
		Kind: GateKindMemoryStore,
		Meta: GatekeeperMeta{Sensitivity: "secret"},
	}

	decision, err := gk.Check(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if decision.Decision != DecisionAllow {
		t.Errorf("expected allow, got: %s", decision.Decision)
	}
}

func TestDenyAllGatekeeper(t *testing.T) {
	gk := &DenyAllGatekeeper{}

	req := GatekeeperCheckRequest{
		Kind: GateKindMemoryStore,
		Meta: GatekeeperMeta{Sensitivity: "public"},
	}

	decision, err := gk.Check(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if decision.Decision != DecisionDeny {
		t.Errorf("expected deny, got: %s", decision.Decision)
	}
}