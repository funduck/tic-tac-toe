// @title		Tic Tac Toe API
// @version		1.0
// @description	Multiplayer Tic Tac Toe game server
// @host		localhost:8080
// @BasePath	/
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Services and handlers
	gameRepo := game.NewMemoryRepo()
	gameSvc := game.NewGameService(gameRepo, logger)
	gameHandler := server.NewGameHandler(gameSvc, logger)

	tokenService := auth.NewAccessTokenService("secret", "tic-tac-toe")
	userRepo := user.NewMemoryUserRepo()
	userSvc := user.NewUserService(userRepo, tokenService, logger)
	userHandler := server.NewUserHandler(userSvc, logger)

	authMiddleware := server.AuthMiddleware(tokenService)

	// Router setup
	router := newRouter()
	router.Route("/api", func(r chi.Router) {
		server.GameRouter(r, gameHandler, authMiddleware)
		server.UserRouter(r, userHandler)
	})

	// Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		logger.Info("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
		}
	}()

	// Graceful shutdown
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGTERM, syscall.SIGINT)

	<-shutdownChannel
	logger.Info("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Error("Server shutdown error", "error", err)
	} else {
		logger.Info("Server gracefully stopped")
	}
}

// Basic router setup with Swagger UI and OpenAPI spec endpoint
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

	// Health check endpoint (dummy for now)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}
