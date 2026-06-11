package app

import (
	"context"
	"fmt"
	"os"

	"github.com/funduck/tic-tac-toe/cmd/client/lib"
	"github.com/funduck/tic-tac-toe/internal/game"
)

// CreateOrJoinGame handles game creation or joining
func CreateOrJoinGame(ctx context.Context, gameSvc *lib.GameService, displaySvc *lib.DisplayService) error {
	if lib.GameID == "" {
		// Try join any game first
		g, err := gameSvc.JoinAnyGame(ctx, lib.UserID)
		if err == nil {
			lib.GameState = g
			displaySvc.PrintInfo(fmt.Sprintf("Joined game: %s with %s", g.GetID(), g.GetOpponentID()))
			return nil
		}

		// Create new game
		g, err = gameSvc.CreateGame(ctx, lib.UserID, lib.Private)
		if err != nil {
			return fmt.Errorf("failed to create game: %w", err)
		}
		lib.GameState = g
		displaySvc.PrintInfo(fmt.Sprintf("Game created: %s\nShare this ID with your opponent.", g.GetID()))
		return nil
	}

	// Join existing game
	g, err := gameSvc.JoinGame(ctx, lib.GameID, lib.UserID)
	if err != nil {
		return fmt.Errorf("failed to join or get game: %w", err)
	}
	lib.GameState = g
	displaySvc.PrintInfo(fmt.Sprintf("Joined game: %s with %s", g.GetID(), g.GetOpponentID()))

	return nil
}

func WaitForOpponent(ctx context.Context, gameSvc *lib.GameService, displaySvc *lib.DisplayService) error {
	g := lib.GameState
	if g.GetStatus() == game.StatusWaiting { // TODO safe types for enums
		displaySvc.PrintInfo("⏳ Waiting for opponent to join...")
		g, err := gameSvc.PollUntil(ctx, g.GetID(), func(g *lib.Game) bool {
			return g.GetStatus() == game.StatusInProgress
		}, 5)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: failed to wait for opponent: %v\n", err)
			os.Exit(1)
		}
		lib.GameState = g
	}
	return nil
}
