package app

import (
	"context"
	"fmt"

	"github.com/funduck/tic-tac-toe/cmd/client/lib"
)

// CreateOrJoinGame handles game creation or joining
func CreateOrJoinGame(ctx context.Context, gameSvc *lib.GameService, displaySvc *lib.DisplayService) error {
	if lib.GameID == "" {
		// Create new game
		g, err := gameSvc.CreateGame(ctx, lib.UserID)
		if err != nil {
			return fmt.Errorf("failed to create game: %w", err)
		}
		lib.GameState = g
		displaySvc.PrintInfo(fmt.Sprintf("Game created: %s\nShare this ID with your opponent.", g.GetID()))

		if g.GetStatus() == "waiting" { // TODO safe types for enums
			displaySvc.PrintInfo("⏳ Waiting for opponent to join...")
			g, err = gameSvc.PollUntil(ctx, g.GetID(), func(g *lib.Game) bool {
				return g.GetStatus() != "waiting"
			}, 5)
			if err != nil {
				return fmt.Errorf("failed to wait for opponent: %w", err)
			}
			lib.GameState = g
		}
		return nil
	}

	// Join existing game
	g, err := gameSvc.JoinGame(ctx, lib.GameID, lib.UserID)
	if err != nil {
		return fmt.Errorf("failed to join or get game: %w", err)
	}
	lib.GameState = g

	return nil
}
