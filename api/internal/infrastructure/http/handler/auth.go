// Package handler - HTTP handlers for auth endpoints.
package handler

import (
	"encoding/json"
	"net/http"

	appauth "github.com/v4lss/animas/internal/application/auth"
	"github.com/v4lss/animas/internal/domain/user"
	"github.com/v4lss/animas/pkg/jwt"
	"github.com/v4lss/animas/pkg/response"
)

// AuthHandler holds dependencies for auth routes.
type AuthHandler struct {
	userRepo user.Repository
	jwtSvc   *jwt.Service
}

func NewAuthHandler(userRepo user.Repository, jwtSvc *jwt.Service) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtSvc: jwtSvc}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate password length
	if len(body.Password) < 8 {
		response.Error(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	u, err := appauth.Register(r.Context(), h.userRepo, appauth.RegisterInput{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, map[string]string{"id": u.ID, "email": u.Email})
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	out, err := appauth.Login(r.Context(), h.userRepo, h.jwtSvc, appauth.LoginInput{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]any{
		"token": out.Token,
		"user":  map[string]string{"id": out.User.ID, "email": out.User.Email},
	})
}
