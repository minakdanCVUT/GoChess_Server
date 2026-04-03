package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minakdanCVUT/GoChess/internal/apperr"
)

type contextKey string

const UserIDKey contextKey = "user_id"

var jwtSecret = []byte("L8thPONOMLXik2zfkF1SzxrpAqS2c137pJzTqAOBvka/JkYovv+Mnd0wMPFumGoaXSQIypehdQXsr/zqm7hkEi0jYV0hBpp1vrWDQcpwjgesntK1j3NucGdW5I5m1YntS/9VNFXprlJ5+hBKJhSdah14y7OMeS16W7M6PoV2hxedP6Aa4I3+ZVDTrX46mVQaGBpTx4ZGMFId5LCB9HvkbmujM143F3fPqVDAPTpZDcvR7Ad1fv2VBwJQ2dLEG+iRZkNjBomHohF+r9R21N1nfl8DeslGWgTdOAm9wikOTbHZXk4aIbhcKiSvLNccO93ScW7MvC6W63jn60oXtlPvRw==%")

func (s *UserService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *UserService) ExtractUserIDFromContext(ctx context.Context) (pgtype.UUID, error) {
	var zero pgtype.UUID
	var userID pgtype.UUID
	userIdStr, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return zero, apperr.ErrUnauthorized
	}
	if err := userID.Scan(userIdStr); err != nil {
		return zero, apperr.ErrUnauthorized
	}
	return userID, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Достаем заголовок Authorization: Bearer <token>
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			apperr.HandleError(w, apperr.ErrUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Проверяем подпись и срок годности
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			apperr.HandleError(w, apperr.ErrUnauthorized)
			return
		}

		// 3. Вытаскиваем user_id из токена и кладем в контекст
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
