package app

import (
	"owlchat/backend/internal/config"
	"owlchat/backend/internal/modules/auth"
	"owlchat/backend/internal/modules/chat"
	"owlchat/backend/internal/modules/message"
	"owlchat/backend/internal/platform/httpx"

	"github.com/go-chi/chi/v5"
)

type App struct {
	r *chi.Mux
}

func New(cfg config.Config) *App {
	r := chi.NewRouter()
	httpx.UseCommonMiddleware(r)

	authSvc := auth.NewService(cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)
	chatStore := chat.NewMemoryStore()
	chatHandler := chat.NewHandler(chatStore)
	msgStore := message.NewMemoryStore()
	msgHandler := message.NewHandler(msgStore, chatStore)

	r.Get("/healthz", httpx.Healthz)
	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/auth", authHandler.Routes())
		api.Mount("/chats", chatHandler.Routes())
		api.Mount("/messages", msgHandler.Routes())
	})

	return &App{r: r}
}

func (a *App) Router() *chi.Mux { return a.r }
