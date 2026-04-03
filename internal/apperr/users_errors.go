package apperr

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// реализация интерфейса error
func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message + ": " + msg,
	}
}

var (
	ErrUserNotFound         = NewAppError(http.StatusNotFound, "Пользователь не найден")
	ErrInvalidCredentials   = NewAppError(http.StatusUnauthorized, "Неверный логин или пароль")
	ErrInternal             = NewAppError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	ErrJsonParsing          = NewAppError(http.StatusBadRequest, "Ошибка парсинга JSON")
	ErrValidate             = NewAppError(http.StatusBadRequest, "Ошибка валидации")
	ErrEmailOrUsernameInUse = NewAppError(http.StatusConflict, "Email или Username уже заняты")
	ErrUnauthorized         = NewAppError(http.StatusUnauthorized, "Вы не авторизованы")
)

func HandleError(w http.ResponseWriter, err error) {
	var appErr *AppError

	if errors.As(err, &appErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.Code)
		json.NewEncoder(w).Encode(appErr)
		return
	}

	log.Printf("Unhandled error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Something went wrong"))
}
