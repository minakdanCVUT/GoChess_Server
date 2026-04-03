package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minakdanCVUT/GoChess/internal/db"
	"github.com/minakdanCVUT/GoChess/internal/handler"
	"github.com/minakdanCVUT/GoChess/internal/service"
)

func main() {
	ctx := context.Background()
	connStr := "postgres://user:pass@localhost:55432/chess_db?sslmode=disable"

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Failed to parse pool config:", err)
	}

	config.MaxConns = 10
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Failed to create connection pool:", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	queries := db.New(pool)

	usersService := service.NewUserService(queries)

	userHandler := handler.NewUsersHandler(usersService)

	router := handler.RegisterUserRoutes(userHandler)

	log.Println("Server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
