package handler

import (
	"net/http"
)

// RegisterRoutes собирает все обработчики в один роутер (Mux).
// Мы возвращаем http.Handler — это стандартный интерфейс для всех роутеров в Go.
func RegisterUserRoutes(userH *UsersHandler) http.Handler {
	mux := http.NewServeMux()

	// Группируем роуты здесь. Если завтра решишь сменить префикс на /api/v1,
	// ты сделаешь это в одном месте.
	mux.HandleFunc("POST /users/register", userH.CreateUser)
	mux.HandleFunc("POST /users/login", userH.LoginUser)
	mux.HandleFunc("GET /user/profile/{user_id}", userH.GetProfile)
	return mux
}
