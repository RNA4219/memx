package service

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrPolicyDenied    = errors.New("policy denied")
	ErrNeedsHuman      = errors.New("needs human review")
)
