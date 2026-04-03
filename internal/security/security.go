package security

import (
	"context"
	"fmt"
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// taking out the Authorization header: Bearer <token>
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			apperr.HandleError(w, apperr.ErrUnauthorized())
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

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
