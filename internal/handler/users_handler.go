package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/minakdanCVUT/GoChess/internal/apperr"
	"github.com/minakdanCVUT/GoChess/internal/db"
	"github.com/minakdanCVUT/GoChess/internal/handler/requests"
	"github.com/minakdanCVUT/GoChess/internal/handler/responses"
	"github.com/minakdanCVUT/GoChess/internal/security"
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

// CreateUser godoc
// @Summary      Register a new user
// @Description  Create a new user account and get an auth token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body requests.CreateUserRequest true "User registration data"
// @Success      201 {object} responses.AuthResponse
// @Failure      400 {object} apperr.AppError
// @Router       /users/register [post]
func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.HandleError(w, apperr.ErrJsonParsing())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperr.HandleError(w, apperr.ErrValidate().WithMessage(err.Error()))
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

// LoginUser godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body requests.LoginUserRequest true "Login credentials"
// @Success      200 {object} responses.AuthResponse
// @Failure      401 {object} apperr.AppError
// @Router       /users/login [post]
func (h *UsersHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request requests.LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apperr.HandleError(w, apperr.ErrJsonParsing())
		return
	}

	if err := h.validate.Struct(request); err != nil {
		apperr.HandleError(w, apperr.ErrValidate().WithMessage(err.Error()))
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

// GetProfile godoc
// @Summary      Get current user profile
// @Description  Returns profile data for the authenticated user
// @Tags         users
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} responses.UserResponse
// @Failure      401 {object} apperr.AppError
// @Router       /users/profile [get]
func (h *UsersHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// extract user id, that AuthMiddleware put in context from JWT
	userID, err := security.ExtractUserIDFromContext(r.Context())
	if err != nil {
		apperr.HandleError(w, err)
		return
	}

	user, err := h.service.Profile(r.Context(), userID)
	if err != nil {
		apperr.HandleError(w, err)
		return
	}

	response := responses.UserResponse{
		ID:            user.ID.String(),
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
