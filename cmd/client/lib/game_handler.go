package lib

import (
	"bufio"
	"context"
	"fmt"
)

// CreateOrJoinGame handles game creation or joining
func CreateOrJoinGame(ctx context.Context, gameSvc *GameService, displaySvc *DisplayService) error {
	if GameID == "" {
		// Create new game
		g, err := gameSvc.CreateGame(ctx, UserID)
		if err != nil {
			return fmt.Errorf("failed to create game: %w", err)
		}
		GameState = g
		displaySvc.PrintInfo(fmt.Sprintf("Game created: %s\nShare this ID with your opponent.", g.GetID()))

		if g.GetStatus() == "waiting" { // TODO safe types for enums
			displaySvc.PrintInfo("⏳ Waiting for opponent to join...")
			g, err = gameSvc.PollUntil(ctx, g.GetID(), func(g *Game) bool {
				return g.GetStatus() != "waiting"
			}, 5)
			if err != nil {
				return fmt.Errorf("failed to wait for opponent: %w", err)
			}
			GameState = g
		}
		return nil
	}

	// Join existing game
	g, err := gameSvc.JoinGame(ctx, GameID, UserID)
	if err != nil {
		// May already be joined — fetch current state
		g, err = gameSvc.GetGame(ctx, GameID)
		if err != nil {
			return fmt.Errorf("failed to join or get game: %w", err)
		}
		GameState = g
	}
	return nil
}

// PlayGame runs the main game loop
func PlayGame(ctx context.Context, gameSvc *GameService, displaySvc *DisplayService, inputSvc *InputService, scanner *bufio.Scanner) error {
	for GameState.GetStatus() == "in_progress" {
		// Wait for our turn
		if GameState.GetCurrentPlayerID() != UserID {
			displaySvc.PrintStatus()
			var err error
			g, err := gameSvc.PollUntil(ctx, GameState.GetID(), func(g *Game) bool {
				return g.GetStatus() != "in_progress" || g.GetCurrentPlayerID() == UserID
			}, 5)
			if err != nil {
				return fmt.Errorf("failed to poll game state: %w", err)
			}
			GameState = g
			displaySvc.PrintBoard()
			continue
		}

		// Our turn
		displaySvc.PrintStatus()
		fmt.Print("Your move (row col) or 'q' to give up: ")

		row, col, giveUp, err := inputSvc.PromptMove(scanner)
		if err != nil {
			displaySvc.PrintError(err.Error())
			continue
		}

		if giveUp {
			g, err := gameSvc.GiveUp(ctx, GameState.GetID(), UserID)
			if err != nil {
				return fmt.Errorf("failed to give up: %w", err)
			}
			GameState = g
			break
		}

		g, err := gameSvc.MakeMove(ctx, GameState.GetID(), UserID, row, col)
		if err != nil {
			displaySvc.PrintError(fmt.Sprintf("Move rejected: %v", err))
			// Refresh state after rejection
			g, err := gameSvc.GetGame(ctx, GameState.GetID())
			if err != nil {
				return fmt.Errorf("failed to refresh game state: %w", err)
			}
			GameState = g
			continue
		}
		GameState = g
		displaySvc.PrintBoard()
	}

	return nil
}
