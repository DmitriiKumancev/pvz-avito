package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dkumancev/avito-pvz/pkg/application/services"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type UserHandler struct {
	userService services.UserService
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"message":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	var role domain.UserRole
	switch req.Role {
	case "employee":
		role = domain.EmployeeRole
	case "moderator":
		role = domain.ModeratorRole
	default:
		http.Error(w, `{"message":"Неверная роль пользователя"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), req.Email, req.Password, role)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	resp := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  string(user.Role),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"message":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	token, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"message":"Неверные учетные данные"}`, http.StatusUnauthorized)
		return
	}

	resp := TokenResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"message":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	var req DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	var role domain.UserRole
	switch req.Role {
	case "employee":
		role = domain.EmployeeRole
	case "moderator":
		role = domain.ModeratorRole
	default:
		http.Error(w, `{"message":"Неверная роль пользователя"}`, http.StatusBadRequest)
		return
	}

	token, err := h.userService.DummyLogin(r.Context(), role)
	if err != nil {
		http.Error(w, `{"message":"Ошибка при создании токена"}`, http.StatusInternalServerError)
		return
	}

	resp := TokenResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
