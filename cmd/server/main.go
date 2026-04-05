package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/minakdanCVUT/GoChess/internal/db"
	"github.com/minakdanCVUT/GoChess/internal/handler"
	"github.com/minakdanCVUT/GoChess/internal/security"
	"github.com/minakdanCVUT/GoChess/internal/service"
	"github.com/minakdanCVUT/GoChess/internal/socket"
)

// @title           GoChess API
// @version         1.0
// @description     Chess Game Server with WebSockets and JWT Auth.
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                Type 'Bearer ' followed by your JWT token
func main() {
	godotenv.Load()
	security.Init()

	ctx := context.Background()
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to parse pool config:", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Failed to create connection pool:", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	queries := db.New(pool)

	gamesService := service.NewGameService(queries)

	hub := socket.NewHub(gamesService)
	go hub.Run()

	config.MaxConns = 10
	config.MaxConnIdleTime = 5 * time.Minute

	usersService := service.NewUserService(queries)

	userHandler := handler.NewUsersHandler(usersService)

	router := handler.RegisterRoutes(userHandler, hub)

	addr := os.Getenv("SERVER_ADDR")
	log.Printf("Server started on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
