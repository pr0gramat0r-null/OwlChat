package httpx

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func UseCommonMiddleware(r interface{ Use(...func(http.Handler) http.Handler) }) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func Healthz(w http.ResponseWriter, _ *http.Request) {
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
