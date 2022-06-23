package handler

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"time"
)

func MakeRouter(ctx context.Context, handler *Handler) http.Handler {
	var router = chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second)) //nolint:gomnd
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{http.MethodGet},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		//nolint:gomnd
		MaxAge: 300, // Maximum value not ignored by any of major browsers
	}))

	setBasic(router, handler)

	router.Get("/encrypt", handler.Encrypt)
	router.Get("/decrypt", handler.Decrypt)

	return router
}

func setBasic(r *chi.Mux, handler *Handler) {
	r.Get("/", handler.OK)
	r.NotFound(handler.NotFound)
}
