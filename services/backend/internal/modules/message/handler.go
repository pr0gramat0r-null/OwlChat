package message

import (
	"errors"
	"net/http"
	"strings"

	"owlchat/backend/internal/modules/chat"
	"owlchat/backend/internal/platform/httpx"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store *MemoryStore
	chats *chat.MemoryStore
}

func NewHandler(store *MemoryStore, chats *chat.MemoryStore) *Handler {
	return &Handler{store: store, chats: chats}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.send)
	r.Get("/{chatID}", h.list)
	return r
}

type sendMessageRequest struct {
	ID string `json:"id"`
	ChatID string `json:"chat_id"`
	SenderID string `json:"sender_id"`
	ClientMsgID string `json:"client_msg_id"`
	Body string `json:"body"`
}

func validate(req sendMessageRequest) error {
	if req.ChatID == "" || req.SenderID == "" || req.ClientMsgID == "" { return errors.New("missing required fields") }
	if n := len(strings.TrimSpace(req.Body)); n == 0 || n > 4000 { return errors.New("body size out of bounds") }
	return nil
}

func (h *Handler) send(w http.ResponseWriter, r *http.Request) {
	var req sendMessageRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.JSON(w, http.StatusBadRequest, map[string]string{"error":"invalid json"})
		return
	}
	if err := validate(req); err != nil {
		httpx.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if _, ok := h.chats.Get(req.ChatID); !ok {
		httpx.JSON(w, http.StatusNotFound, map[string]string{"error":"chat not found"})
		return
	}
	msg := h.store.Add(Message(req))
	httpx.JSON(w, http.StatusCreated, msg)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chatID")
	httpx.JSON(w, http.StatusOK, map[string]any{"items": h.store.List(chatID)})
}
