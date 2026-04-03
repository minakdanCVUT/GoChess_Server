package service

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minakdanCVUT/GoChess/internal/apperr"
	"github.com/minakdanCVUT/GoChess/internal/db"
	"github.com/minakdanCVUT/GoChess/internal/security"
	"golang.org/x/crypto/bcrypt"
)

const postgresUniqueErrorCode = "23505"

type UserService struct {
	queries *db.Queries
}

func NewUserService(q *db.Queries) *UserService {
	return &UserService{queries: q}
}

func (s *UserService) Login(ctx context.Context, login string, password string) (*db.User, string, error) {
	user, err := s.queries.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, "", apperr.ErrUserNotFound()
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", apperr.ErrInvalidCredentials()
	}

	token, _ := security.GenerateToken(user.ID.String())

	log.Printf("User logged in, username: %s", user.Username)
	return &user, token, nil
}

func (s *UserService) Register(ctx context.Context, params *db.CreateUserParams) (*db.User, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", apperr.ErrInternal()
	}
	params.Password = string(hashedPassword)

	user, err := s.queries.CreateUser(ctx, *params)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == postgresUniqueErrorCode {
			return nil, "", apperr.ErrEmailOrUsernameInUse()
		}
		return nil, "", apperr.ErrInternal()
	}

	token, _ := security.GenerateToken(user.ID.String())

	return &user, token, nil
}

func (s *UserService) Profile(ctx context.Context, userID pgtype.UUID) (*db.User, error) {
	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrUserNotFound()
		}
		return nil, apperr.ErrInternal()
	}
	return &user, nil
}
