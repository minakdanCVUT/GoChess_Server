package security

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minakdanCVUT/GoChess/internal/apperr"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const GameTypeKey contextKey = "game_type"

type GameType string

const (
	GameTypeBlitz  GameType = "blitz"
	GameTypeRapid  GameType = "rapid"
	GameTypeRating GameType = "rating"
)

func AllGameTypes() []GameType {
	return []GameType{
		GameTypeBlitz,
		GameTypeRapid,
		GameTypeRating,
	}
}

func IsValidGameType(t string) bool {
	switch GameType(t) {
	case GameTypeBlitz, GameTypeRapid, GameTypeRating:
		return true
	}
	return false
}

func GetGameTypeFromString(t string) GameType {
	var zero GameType
	switch t {
	case "blitz":
		return GameTypeBlitz
	case "rapid":
		return GameTypeRapid
	case "rating":
		return GameTypeRating
	}
	return zero
}

var jwtSecret []byte

func Init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	jwtSecret = []byte(secret)
}

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ExtractUserIDFromContext(ctx context.Context) (pgtype.UUID, error) {
	var zero pgtype.UUID
	var userID pgtype.UUID
	userIdStr, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return zero, apperr.ErrUnauthorized()
	}
	if err := userID.Scan(userIdStr); err != nil {
		return zero, apperr.ErrUnauthorized()
	}
	return userID, nil
}

func ExtractGameTypeFromContext(ctx context.Context) (string, error) {
	gameType, ok := ctx.Value(GameTypeKey).(string)
	if !ok {
		return "", apperr.ErrInternal()
	}
	return gameType, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// if headers are empty, searching in the url queries for websocket
		if tokenStr == "" {
			tokenStr = r.URL.Query().Get("token")
		}

		if tokenStr == "" {
			apperr.HandleError(w, apperr.ErrUnauthorized())
			return
		}

		// check the signature and expiration date
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			apperr.HandleError(w, apperr.ErrUnauthorized())
			return
		}

		// taking out user_id from token and put to the context
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func QueriesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gameType := r.URL.Query().Get("game_type")

		if !IsValidGameType(gameType) || gameType == "" {
			apperr.HandleError(w, apperr.ErrInternal())
			return
		}
		ctx := context.WithValue(r.Context(), GameTypeKey, gameType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
