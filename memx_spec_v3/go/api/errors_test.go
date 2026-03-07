package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"memx/service"
)

// TestMapError_ErrorCodeMapping はエラーコードのマッピングをテストする。
// 仕様書: requirements-api.md 6-4. エラーモデル
func TestMapError_ErrorCodeMapping(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode ErrorCode
	}{
		{
			name:     "ErrNotFound maps to NOT_FOUND",
			err:      service.ErrNotFound,
			wantCode: CodeNotFound,
		},
		{
			name:     "ErrInvalidArgument maps to INVALID_ARGUMENT",
			err:      service.ErrInvalidArgument,
			wantCode: CodeInvalidArgument,
		},
		{
			name:     "ErrPolicyDenied maps to GATEKEEP_DENY",
			err:      service.ErrPolicyDenied,
			wantCode: CodeGatekeepDeny,
		},
		{
			name:     "ErrNeedsHuman maps to GATEKEEP_DENY",
			err:      service.ErrNeedsHuman,
			wantCode: CodeGatekeepDeny,
		},
		{
			name:     "ErrFeatureDisabled maps to INTERNAL",
			err:      service.ErrFeatureDisabled,
			wantCode: CodeInternal,
		},
		{
			name:     "unknown error maps to INTERNAL",
			err:      errors.New("unknown error"),
			wantCode: CodeInternal,
		},
		{
			name:     "error containing 'invalid argument' maps to INVALID_ARGUMENT",
			err:      errors.New("something: invalid argument: detail"),
			wantCode: CodeInvalidArgument,
		},
		{
			name:     "nil error returns nil",
			err:      nil,
			wantCode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapError(tt.err)
			if tt.err == nil {
				if result != nil {
					t.Error("expected nil for nil error")
				}
				return
			}
			if result == nil {
				t.Fatal("expected non-nil error")
			}
			if result.Code != tt.wantCode {
				t.Errorf("expected code %q, got: %q", tt.wantCode, result.Code)
			}
		})
	}
}

// TestHTTP_ErrorStatusCodes はエラーコードからHTTPステータスへの変換をテストする。
// 仕様書: requirements-api.md 6-4. エラーモデル
// - INVALID_ARGUMENT → 400
// - NOT_FOUND → 404
// - CONFLICT → 409
// - GATEKEEP_DENY → 403
// - INTERNAL → 500
func TestHTTP_ErrorStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		code       ErrorCode
		wantStatus int
	}{
		{
			name:       "INVALID_ARGUMENT maps to 400",
			code:       CodeInvalidArgument,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NOT_FOUND maps to 404",
			code:       CodeNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "CONFLICT maps to 409",
			code:       CodeConflict,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "GATEKEEP_DENY maps to 403",
			code:       CodeGatekeepDeny,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "INTERNAL maps to 500",
			code:       CodeInternal,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// HTTPサーバーのwriteErr関数を使用してステータスを確認
			rec := httptest.NewRecorder()
			writeErr(rec, &Error{Code: tt.code, Message: "test error"})

			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got: %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

// TestHTTP_ErrorResponseFormat はエラーレスポンスのフォーマットをテストする。
// 仕様書: requirements-api.md 6-4. エラーモデル
// 共通フォーマット: {"code":"...","message":"...","details":{}}
func TestHTTP_ErrorResponseFormat(t *testing.T) {
	rec := httptest.NewRecorder()
	writeErr(rec, &Error{
		Code:    CodeInvalidArgument,
		Message: "title is required",
		Details: map[string]interface{}{"field": "title"},
	})

	// JSONフォーマット確認
	var errResp Error
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if errResp.Code != CodeInvalidArgument {
		t.Errorf("expected code %q, got: %q", CodeInvalidArgument, errResp.Code)
	}
	if errResp.Message != "title is required" {
		t.Errorf("expected message 'title is required', got: %q", errResp.Message)
	}
}

// TestHTTP_NilErrorResponse はnilエラーの場合の挙動をテストする。
func TestHTTP_NilErrorResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	writeErr(rec, nil)

	// nilの場合でも500が返る
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got: %d", rec.Code)
	}
}