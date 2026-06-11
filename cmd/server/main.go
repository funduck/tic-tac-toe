// @title		Tic Tac Toe API
// @version		1.0
// @description	Multiplayer Tic Tac Toe game server
// @host		localhost:8080
// @BasePath	/
package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/funduck/tic-tac-toe/internal/auth"
	"github.com/funduck/tic-tac-toe/internal/game"
	"github.com/funduck/tic-tac-toe/internal/server"
	"github.com/funduck/tic-tac-toe/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	gameRepo := game.NewMemoryRepo()
	gameSvc := game.NewGameService(gameRepo, logger)
	gameHandler := server.NewGameHandler(gameSvc, logger)

	tokenService := auth.NewAccessTokenService("secret", "tic-tac-toe")
	userRepo := user.NewMemoryUserRepo()
	userSvc := user.NewUserService(userRepo, tokenService) // TODO add logger
	userHandler := server.NewUserHandler(userSvc, logger)

	router := newRouter()

	router.Route("/api", func(r chi.Router) {
		server.GameRouter(r, gameHandler)
		server.UserRouter(r, userHandler)
	})

	logger.Info("Server listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("Server error", "error", err)
	}
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

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
