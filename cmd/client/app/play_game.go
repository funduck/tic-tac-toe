package app

import (
	"context"
	"fmt"

	"github.com/funduck/tic-tac-toe/client/lib"
)

// PlayGame runs the main game loop
func PlayGame(ctx context.Context, gameSvc *lib.GameService, displaySvc *lib.DisplayService, inputSvc *lib.InputService) error {
	for lib.GameState.GetStatus() == "in_progress" {
		// Wait for our turn
		if lib.GameState.GetCurrentPlayerID() != lib.UserID {
			displaySvc.PrintStatus()
			var err error
			g, err := gameSvc.PollUntil(ctx, lib.GameState.GetID(), func(g *lib.Game) bool {
				return g.GetStatus() != "in_progress" || g.GetCurrentPlayerID() == lib.UserID
			}, 5)
			if err != nil {
				return fmt.Errorf("failed to poll game state: %w", err)
			}
			lib.GameState = g
			displaySvc.PrintBoard()
			continue
		}

		// Our turn
		displaySvc.PrintStatus()
		fmt.Print("Your move (row col) or 'q' to give up: ")

		row, col, giveUp, err := inputSvc.PromptMove()
		if err != nil {
			displaySvc.PrintError(err.Error())
			continue
		}

		if giveUp {
			g, err := gameSvc.GiveUp(ctx, lib.GameState.GetID(), lib.UserID)
			if err != nil {
				return fmt.Errorf("failed to give up: %w", err)
			}
			lib.GameState = g
			break
		}

		g, err := gameSvc.MakeMove(ctx, lib.GameState.GetID(), lib.UserID, row, col)
		if err != nil {
			displaySvc.PrintError(fmt.Sprintf("Move rejected: %v", err))
			// Refresh state after rejection
			g, err := gameSvc.GetGame(ctx, lib.GameState.GetID())
			if err != nil {
				return fmt.Errorf("failed to refresh game state: %w", err)
			}
			lib.GameState = g
			continue
		}
		lib.GameState = g
		displaySvc.PrintBoard()
	}

	return nil
}
