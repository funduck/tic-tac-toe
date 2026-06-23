package app

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/funduck/tic-tac-toe/client/lib"
)

func Start(ctx context.Context, gameSvc *lib.GameService, userSvc *lib.UserService, displaySvc *lib.DisplayService, inputSvc *lib.InputService) {
	// Authenticate user (login or signup)
	if err := LoginOrSignup(ctx, userSvc, displaySvc); err != nil {
		fatal("Authentication failed", err)
	}

	// Create or join game
	if err := CreateOrJoinGame(ctx, gameSvc, displaySvc); err != nil {
		fatal("Error", err)
	}

	// Read stdin once and share the channel: both the waiting phase and the game
	// loop need to react to 'q' while concurrently polling the server.
	input := inputSvc.Lines()

	// Wait for opponent if needed
	if err := WaitForOpponent(ctx, gameSvc, displaySvc, input); err != nil {
		if errors.Is(err, ErrUserQuit) {
			displaySvc.PrintInfo("👋 Left the game.")
			return
		}
		fatal("Error", err)
	}

	// Game loop (renders the board and reads moves)
	if err := PlayGame(ctx, gameSvc, displaySvc, input); err != nil {
		fatal("Error", err)
	}

	// Final board + result
	displaySvc.RenderFrame()
	displaySvc.PrintResult()
}

// fatal prints a user-facing error and exits.
func fatal(prefix string, err error) {
	fmt.Fprintf(os.Stderr, "❌ %s: %s\n", prefix, lib.FriendlyMessage(err))
	os.Exit(1)
}
