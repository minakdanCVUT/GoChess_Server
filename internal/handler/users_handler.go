package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/go-playground/validator/v10"
	"github.com/minakdanCVUT/GoChess/internal/apperr"
	"github.com/minakdanCVUT/GoChess/internal/db"
	"github.com/minakdanCVUT/GoChess/internal/handler/requests"
	"github.com/minakdanCVUT/GoChess/internal/handler/responses"
	"github.com/minakdanCVUT/GoChess/internal/service"
)

type UsersHandler struct {
	service  *service.UserService
	validate *validator.Validate
}

func NewUsersHandler(s *service.UserService) *UsersHandler {
	return &UsersHandler{
		service:  s,
		validate: validator.New(),
	}
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.HandleError(w, apperr.ErrJsonParsing)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperr.HandleError(w, apperr.ErrValidate.WithMessage(err.Error()))
		return
	}

	params := db.CreateUserParams{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	user, token, err := h.service.Register(r.Context(), &params)
	if err != nil {
		apperr.HandleError(w, err)
		return
	}

	response := responses.AuthResponse{
		Token:    token,
		UserId:   user.ID.String(),
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UsersHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request requests.LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperr.HandleError(w, apperr.ErrJsonParsing)
		return
	}

	if err := h.validate.Struct(request); err != nil {
		apperr.HandleError(w, apperr.ErrValidate.WithMessage(err.Error()))
		return
	}

	user, token, err := h.service.Login(r.Context(), request.Login, request.Password)
	if err != nil {
		apperr.HandleError(w, err)
		return
	}

	response := responses.AuthResponse{
		Token:    token,
		UserId:   user.ID.String(),
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UsersHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	var userID pgtype.UUID

	err := userID.Scan(r.PathValue("user_id"))
	if err != nil {
		apperr.HandleError(w, apperr.ErrUserNotFound)
		return
	}

	user, err := h.service.Profile(r.Context(), userID)
	if err != nil {
		apperr.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
