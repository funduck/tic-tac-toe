package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/funduck/tic-tac-toe/docs"
)

// NewRouter builds and returns the chi router with all routes mounted.
func NewRouter(h *GameHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		// here we will inject auth middleware in the future if needed

		r.Use(loggerMiddleware)

		r.Route("/api", func(r chi.Router) {
			r.Post("/games", h.CreateGame)
			r.Post("/games/{gameID}/join", h.JoinGame)
			r.Post("/games/{gameID}/move", h.MakeMove)
			r.Post("/games/{gameID}/giveup", h.GiveUp)
			r.Get("/games/{gameID}", h.GetGame)
		})
	})

	// Serve raw OpenAPI spec (consumed by codegen_client.sh)
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return r
}
