package handler

import (
	"net/http"

	"github.com/minakdanCVUT/GoChess/internal/security"
)

func RegisterUserRoutes(userH *UsersHandler) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /users/register", userH.CreateUser)
	mux.HandleFunc("POST /users/login", userH.LoginUser)

	// Profile is accessible only to authenticated users
	profileHandler := security.AuthMiddleware(http.HandlerFunc(userH.GetProfile))
	mux.Handle("GET /users", profileHandler)

	return mux
}
