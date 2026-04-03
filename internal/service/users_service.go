package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minakdanCVUT/GoChess/internal/apperr"
	"github.com/minakdanCVUT/GoChess/internal/db"
)

const POSTGRES_UNIQUE_ERROR_CODE = "23505"

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

func (s *UserService) Register(ctx context.Context, params *db.CreateUserParams) (*db.User, string, error) {
	user, err := s.queries.CreateUser(ctx, *params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == POSTGRES_UNIQUE_ERROR_CODE {
			return nil, "", apperr.ErrEmailOrUsernameInUse
		}
		return nil, "", apperr.ErrInternal
	}

	token := generateRandomToken(16)

	return &user, token, nil
}

func (s *UserService) Profile(ctx context.Context, userID pgtype.UUID) (*db.User, error) {
	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrUserNotFound
		}
		return nil, apperr.ErrInternal
	}
	return &user, nil
}
