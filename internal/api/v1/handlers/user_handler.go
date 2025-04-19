package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dkumancev/avito-pvz/internal/api/response"
	"github.com/dkumancev/avito-pvz/pkg/application/services"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type RegisterRequest struct {
	Email    string          `json:"email"`
	Password string          `json:"password"`
	Role     domain.UserRole `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DummyLoginRequest struct {
	Role domain.UserRole `json:"role"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if req.Role != domain.EmployeeRole && req.Role != domain.ModeratorRole {
		response.Error(w, http.StatusBadRequest, "Invalid role")
		return
	}

	user, err := h.userService.Register(r.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			response.Error(w, http.StatusBadRequest, "User with this email already exists")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	token, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to authenticate user")
		return
	}

	response.JSON(w, http.StatusOK, token)
}

func (h *UserHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Role != domain.EmployeeRole && req.Role != domain.ModeratorRole {
		response.Error(w, http.StatusBadRequest, "Invalid role")
		return
	}

	token, err := h.userService.DummyLogin(r.Context(), req.Role)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate dummy token")
		return
	}

	response.JSON(w, http.StatusOK, token)
}
