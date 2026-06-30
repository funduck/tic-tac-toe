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

	gameSvc := game.NewGameService(gameRepo, logger, getAfkTimeout(logger))
	gameHandler := server.NewGameHandler(gameSvc, logger)

	tokenService := auth.NewAccessTokenService(getSecret(logger), "tic-tac-toe")

	userRepo := user.NewMemoryUserRepo()
	userSvc := user.NewUserService(userRepo, tokenService, logger)
	userHandler := server.NewUserHandler(userSvc, logger, isSecureCookieEnabled())

	authMiddleware := server.AuthMiddleware(tokenService)

	// Router setup
	router := newRouter(logger)
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
func newRouter(logger *slog.Logger) *chi.Mux {
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
		_, err := w.Write([]byte("OK"))
		if err != nil {
			logger.Error("Failed to write health check response", "error", err)
		}
	})

	return r
}

// getAfkTimeout retrieves the AFK timeout duration from the environment variable "AFK_TIMEOUT".
func getAfkTimeout(logger *slog.Logger) time.Duration {
	afkTimeout := os.Getenv("AFK_TIMEOUT")
	if afkTimeout == "" {
		logger.Warn("AFK_TIMEOUT not set, using default 10s (not recommended for production)")
		return 10 * time.Second
	}
	duration, err := time.ParseDuration(afkTimeout)
	if err != nil {
		logger.Error("Invalid AFK_TIMEOUT value, using default 10s", "error", err)
		return 10 * time.Second
	}
	return duration
}

// getSecret retrieves the JWT secret from the environment variable "JWT_SECRET".
func getSecret(logger *slog.Logger) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logger.Warn("JWT_SECRET not set, using default secret (not recommended for production)")
		return "default-secret"
	}
	return secret
}

// secureCookie determines whether cookies should be marked as Secure (HTTPS only).
// In production, this should be true. In development, it can be false for local testing.
func isSecureCookieEnabled() bool {
	secureCookie := os.Getenv("SECURE_COOKIE")
	return secureCookie == "true"
}
