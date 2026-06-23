package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "github.com/funduck/tic-tac-toe/docs"
)

// UserRouter builds and returns the chi router with all routes mounted.
// Optional middlewares can be passed to be applied before the logger middleware.
// This allows injecting auth or other middlewares while keeping tests simple.
func UserRouter(r chi.Router, h *UserHandler, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		// Apply custom middlewares (e.g., auth in production, none in tests)
		for _, mw := range middlewares {
			r.Use(mw)
		}

		r.Use(LoggerMiddleware)

		r.Post("/users/signup", h.Signup)
		r.Post("/users/login", h.Login)
		r.Post("/users/refresh-token", h.RefreshToken)
	})
}
