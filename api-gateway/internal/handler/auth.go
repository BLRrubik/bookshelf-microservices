package handler

import (
	"net/http"

	"bookshelf/api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	proxy *proxy.ServiceProxy
}

func NewAuthHandler(p *proxy.ServiceProxy) *AuthHandler {
	return &AuthHandler{proxy: p}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyAuthPath(w, r, "/api/v1/auth/register")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyAuthPath(w, r, "/api/v1/auth/login")
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyAuthPath(w, r, "/api/v1/auth/refresh")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyAuthPath(w, r, "/api/v1/auth/logout")
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyAuthPath(w, r, "/api/v1/users/me")
}

func (h *AuthHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyAuthPath(w, r, "/api/v1/users/me")
}

func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	h.proxy.ProxyAuthPath(w, r, "/api/v1/users/"+userID)
}
