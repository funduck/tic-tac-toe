package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "github.com/funduck/tic-tac-toe/docs"
)

// NewRouter builds and returns the chi router with all routes mounted.
// Optional middlewares can be passed to be applied before the logger middleware.
// This allows injecting auth or other middlewares while keeping tests simple.
func NewRouter(h *GameHandler, middlewares ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		// Apply custom middlewares (e.g., auth in production, none in tests)
		for _, mw := range middlewares {
			r.Use(mw)
		}

		r.Use(loggerMiddleware)

		r.Route("/api", func(r chi.Router) {
			r.Post("/games", h.CreateGame)
			r.Post("/games/{gameID}/join", h.JoinGame)
			r.Post("/games/{gameID}/move", h.MakeMove)
			r.Post("/games/{gameID}/giveup", h.GiveUp)
			r.Get("/games/{gameID}", h.GetGame)
		})
	})

	return r
}
