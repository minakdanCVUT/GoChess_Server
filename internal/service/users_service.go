package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/minakdanCVUT/GoChess/internal/apperr"
	"github.com/minakdanCVUT/GoChess/internal/db"
)

type UserService struct {
	queries *db.Queries
}

func NewUserService(q *db.Queries) *UserService {
	return &UserService{queries: q}
}

func generateRandomToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (s *UserService) Login(ctx context.Context, login string, password string) (*db.User, string, error) {
	user, err := s.queries.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, "", apperr.ErrUserNotFound
	}

	if user.Password != password {
		return nil, "", apperr.ErrInvalidCredentials
	}

	token := generateRandomToken(16)

	log.Printf("Залогинился юзер, username - %s", user.Username)
	return &user, token, nil
}
