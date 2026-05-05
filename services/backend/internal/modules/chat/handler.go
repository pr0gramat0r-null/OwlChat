package chat

import (
	"net/http"

	"owlchat/backend/internal/platform/httpx"

	"github.com/go-chi/chi/v5"
)

type Handler struct { store *MemoryStore }

func NewHandler(store *MemoryStore) *Handler { return &Handler{store: store} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.create)
	return r
}

type createChatRequest struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Members []string `json:"members"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createChatRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	chat, err := h.store.Create(Chat{ID: req.ID, Title: req.Title, Members: req.Members})
	if err != nil {
		httpx.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	httpx.JSON(w, http.StatusCreated, chat)
}
