package handler

import (
	"net/http"

	"github.com/minakdanCVUT/GoChess/internal/service"
)

func RegisterUserRoutes(userH *UsersHandler) http.Handler {
	mux := http.NewServeMux()

	// Публичные роуты
	mux.HandleFunc("POST /users/register", userH.CreateUser)
	mux.HandleFunc("POST /users/login", userH.LoginUser)

	// Профиль может смотреть только авторизованный юзер
	profileHandler := service.AuthMiddleware(http.HandlerFunc(userH.GetProfile))
	mux.Handle("GET /users", profileHandler)

	return mux
}
