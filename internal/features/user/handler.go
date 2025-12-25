package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// AUTH REQUIRED
	r.Get("/", h.list)
	r.Get("/{id}", h.getByID)

	// ADMIN ONLY
	r.Post("/", h.create)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)

	return r
}

type CreateRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"123456"`
	FullName string `json:"full_name" example:"Admin User"`
	RoleID   int64  `json:"role_id" example:"1"`
}

// @Summary List users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} User
// @Failure 500 {string} string "internal error"
// @Router /users [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		slog.Error("list users failed", slog.Any("err", err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	for _, u := range users {
		u.Password = ""
	}

	respondJSON(w, http.StatusOK, users)
}

// @Summary Get user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {string} string "not found"
// @Router /users/{id} [get]
func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)

	u, err := h.service.GetUser(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	u.Password = ""
	respondJSON(w, http.StatusOK, u)
}

// @Summary Create user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body createRequest true "User data"
// @Success 201
// @Failure 400 {string} string "bad request"
// @Router /users [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	u := &User{
		Username: req.Username,
		FullName: req.FullName,
		RoleID:   req.RoleID,
	}

	if err := h.service.CreateUser(u, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Update user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param input body createRequest true "User data"
// @Success 200
// @Failure 400 {string} string "bad request"
// @Router /users/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	u := &User{
		ID:       id,
		Username: req.Username,
		FullName: req.FullName,
		RoleID:   req.RoleID,
	}

	if err := h.service.UpdateUser(u, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete user
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {string} string "not found"
// @Router /users/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)

	if err := h.service.DeleteUser(id); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func paramID(r *http.Request) int64 {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}

func respondJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}
