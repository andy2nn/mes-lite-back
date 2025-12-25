package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	auth *AuthService
}

func NewAuthHandler(auth *AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.login)
	r.Post("/refresh", h.refresh)
	return r
}

type loginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"123456"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"uuid-refresh-token"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

// Login
// @Summary Login
// @Description Authenticate user and get access + refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body loginRequest true "Credentials"
// @Success 200 {object} tokenResponse
// @Failure 400 {string} string "bad request"
// @Failure 401 {string} string "invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	access, refresh, user, err := h.auth.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	respondJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		User:         user,
	})
}

// Refresh token
// @Summary Refresh access token
// @Description Refresh JWT using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body refreshRequest true "Refresh token"
// @Success 200 {object} tokenResponse
// @Failure 400 {string} string "bad request"
// @Failure 401 {string} string "invalid refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	access, refresh, err := h.auth.Refresh(req.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	respondJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
