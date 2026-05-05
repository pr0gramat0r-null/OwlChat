package auth

import (
	"net/http"

	"owlchat/backend/internal/platform/httpx"

	"github.com/go-chi/chi/v5"
)

type Handler struct { svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/dev-login", h.devLogin)
	return r
}

type devLoginRequest struct { UserID string `json:"user_id"` }

func (h *Handler) devLogin(w http.ResponseWriter, r *http.Request) {
	var req devLoginRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	t, err := h.svc.IssueAccessToken(req.UserID)
	if err != nil {
		httpx.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"access_token": t})
}
