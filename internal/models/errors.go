package models

import "fmt"

type ErrorCode string

const (
	ErrorCodeTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrorCodePullRExists ErrorCode = "PR_EXISTS"
	ErrorCodePRMerged    ErrorCode = "PR_MERGED"
	ErrorCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrorCodeNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrorCodeNotFound    ErrorCode = "NOT_FOUND"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewAppError(code ErrorCode, message string) AppError {
	return AppError{Code: code, Message: message}
}

var (
	ErrTeamExists  = NewAppError(ErrorCodeTeamExists, "team_name already exists")
	ErrPRExists    = NewAppError(ErrorCodePullRExists, "PR id already exists")
	ErrPRMerged    = NewAppError(ErrorCodePRMerged, "cannot reassign on merged PR")
	ErrNotAssigned = NewAppError(ErrorCodeNotAssigned, "reviewer is not assigned to this PR")
	ErrNoCandidate = NewAppError(ErrorCodeNoCandidate, "no active replacement candidate in team")
	ErrNotFound    = NewAppError(ErrorCodeNotFound, "resource not found")
)
