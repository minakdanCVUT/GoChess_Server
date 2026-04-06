package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/minakdanCVUT/GoChess/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/minakdanCVUT/GoChess/internal/security"
	"github.com/minakdanCVUT/GoChess/internal/socket"
)

func RegisterRoutes(userH *UsersHandler, hub *socket.Hub) http.Handler {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Post("/register", userH.CreateUser)
		r.Post("/login", userH.LoginUser)

		r.Group(func(r chi.Router) {
			r.Use(security.AuthMiddleware)
			r.Get("/profile", userH.GetProfile)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(security.AuthMiddleware)
		r.Use(security.QueriesMiddleware)
		r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			socket.ServeWs(hub, w, r)
		})
	})

	return r
}
