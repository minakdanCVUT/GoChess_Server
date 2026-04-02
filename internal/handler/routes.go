package handler

import (
	"net/http"
)

func RegisterUserRoutes(userH *UsersHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users/register", userH.CreateUser)
	mux.HandleFunc("POST /users/login", userH.LoginUser)
	mux.HandleFunc("GET /user/profile/{user_id}", userH.GetProfile)
	return mux
}
