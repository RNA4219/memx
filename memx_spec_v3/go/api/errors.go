package api

import (
	"errors"
	"strings"

	"memx/service"
)

func mapError(err error) *Error {
	if err == nil {
		return nil
	}
	if errors.Is(err, service.ErrNotFound) {
		return &Error{Code: CodeNotFound, Message: "not found"}
	}
	if errors.Is(err, service.ErrInvalidArgument) {
		return &Error{Code: CodeInvalidArgument, Message: "invalid argument"}
	}
	if errors.Is(err, service.ErrPolicyDenied) {
		return &Error{Code: CodeGatekeepDeny, Message: err.Error()}
	}
	if errors.Is(err, service.ErrNeedsHuman) {
		return &Error{Code: CodeGatekeepDeny, Message: err.Error()}
	}
	if errors.Is(err, service.ErrConflict) {
		return &Error{Code: CodeConflict, Message: err.Error()}
	}
	if errors.Is(err, service.ErrFeatureDisabled) {
		return &Error{Code: CodeInternal, Message: "feature disabled"}
	}
	if strings.Contains(err.Error(), "invalid argument") {
		return &Error{Code: CodeInvalidArgument, Message: err.Error()}
	}
	return &Error{Code: CodeInternal, Message: err.Error()}
}
