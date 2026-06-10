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

	"github.com/funduck/tic-tac-toe/internal/game"
	server "github.com/funduck/tic-tac-toe/internal/http"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	repo := game.NewMemoryRepo()
	svc := game.NewGameService(repo, logger)
	handler := server.NewGameHandler(svc, logger)
	router := server.NewRouter(handler)

	logger.Info("Server listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("Server error", "error", err)
	}
}
