package apperr

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ErrKind string

const (
	KindNotFound     ErrKind = "not_found"
	KindUnauthorized ErrKind = "unauthorized"
	KindInternal     ErrKind = "internal"
	KindBadRequest   ErrKind = "bad_request"
	KindConflict     ErrKind = "conflict"
)

type AppError struct {
	Kind    ErrKind `json:"-"`
	Message string  `json:"message"`
	cause   error
}

func (e *AppError) Error() string {
	return e.Message
}
func (e *AppError) Unwrap() error {
	return e.cause
}

func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{
		Kind:    e.Kind,
		Message: e.Message + ": " + msg,
		cause:   e,
	}
}

func ErrUserNotFound() *AppError {
	return &AppError{Kind: KindNotFound, Message: "User not found"}
}

func ErrInvalidCredentials() *AppError {
	return &AppError{Kind: KindUnauthorized, Message: "Invalid login or password"}
}

func ErrInternal() *AppError {
	return &AppError{Kind: KindInternal, Message: "Internal server error"}
}

func ErrJsonParsing() *AppError {
	return &AppError{Kind: KindBadRequest, Message: "Failed to parse JSON"}
}

func ErrValidate() *AppError {
	return &AppError{Kind: KindBadRequest, Message: "Validation error"}
}

func ErrEmailOrUsernameInUse() *AppError {
	return &AppError{Kind: KindConflict, Message: "Email or username is already taken"}
}

func ErrUnauthorized() *AppError {
	return &AppError{Kind: KindUnauthorized, Message: "Unauthorized"}
}

func httpStatus(kind ErrKind) int {
	switch kind {
	case KindNotFound:
		return http.StatusNotFound
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindBadRequest:
		return http.StatusBadRequest
	case KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func HandleError(w http.ResponseWriter, err error) {
	if appErr, ok := errors.AsType[*AppError](err); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus(appErr.Kind))
		json.NewEncoder(w).Encode(appErr)
		return
	}

	log.Printf("Unhandled error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Something went wrong"))
}
