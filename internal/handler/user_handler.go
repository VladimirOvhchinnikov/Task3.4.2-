package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"projetpostgre/internal/domain"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserRepository domain.UserRepository
}

func NewUserHandler(userRepo domain.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepository: userRepo,
	}
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var createUserReq CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := domain.User{
		Username:  createUserReq.Username,
		Email:     createUserReq.Email,
		Password:  createUserReq.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := h.UserRepository.Create(context.Background(), user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "User created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "id")

	user, err := h.UserRepository.GetByID(context.Background(), userID)
	if err != nil {

		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	var updatedUser domain.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.UserRepository.Update(context.Background(), updatedUser); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if err := h.UserRepository.Delete(context.Background(), id); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) ListHandler(w http.ResponseWriter, r *http.Request) {

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	users, err := h.UserRepository.List(context.Background(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
