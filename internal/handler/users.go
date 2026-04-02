package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	//"strconv"
	"github.com/go-playground/validator/v10"
	"github.com/minakdanCVUT/GoChess/internal/db"
	"github.com/minakdanCVUT/GoChess/internal/handler/requests"
	"github.com/minakdanCVUT/GoChess/internal/handler/responses"
)

type UsersHandler struct {
	queries  *db.Queries
	validate *validator.Validate
}

func NewUsersHandler(q *db.Queries) *UsersHandler {
	return &UsersHandler{
		queries:  q,
		validate: validator.New(),
	}
}

func generateRandomToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "Ошибка валидации: "+err.Error(), http.StatusBadRequest)
		return
	}

	params := db.CreateUserParams{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	log.Printf("Регается пользователь - %s %s.\nUsername - %s, email - %s", params.FirstName, params.LastName, params.Username, params.Email)

	user, err := h.queries.CreateUser(r.Context(), params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Printf("Username или Email уже заняты - %s:%s", params.Username, params.Email)
			http.Error(w, "Username или Email уже заняты", http.StatusConflict)
			return
		}
		http.Error(w, "Не удалось создать пользователя", http.StatusInternalServerError)
		return
	}

	token := generateRandomToken(16)

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
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(request); err != nil {
		http.Error(w, "Ошибка валидации"+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByLogin(r.Context(), request.Login)
	if err != nil {
		http.Error(w, "Не найден пользователем с таким логином", http.StatusNotFound)
		return
	}

	if user.Password != request.Password {
		http.Error(w, "Неверный пароль", http.StatusBadRequest)
		return
	}

	log.Printf("Залогинился юзер, username - %s", user.Username)

	token := generateRandomToken(16)

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
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}
