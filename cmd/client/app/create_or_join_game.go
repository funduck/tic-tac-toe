package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/funduck/tic-tac-toe/client/lib"
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

		// Only fall through to create if no game was available; surface real errors.
		var apiErr *lib.APIError
		if !errors.As(err, &apiErr) || !apiErr.HasCode(openapi.CodeGameNotFound) {
			return fmt.Errorf("failed to join any game: %w", err)
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

// ErrUserQuit signals that the player chose to leave a still-waiting game with
// 'q'. Callers treat it as a clean exit rather than an error.
var ErrUserQuit = errors.New("user quit")

// WaitForOpponent blocks until an opponent joins (status becomes in_progress),
// polling the server while concurrently reading stdin so the player can leave
// with 'q'. Leaving cancels the game via the quit endpoint and returns
// ErrUserQuit.
func WaitForOpponent(ctx context.Context, gameSvc *lib.GameService, displaySvc *lib.DisplayService, input <-chan string) error {
	if lib.GameState.GetStatus() != openapi.StatusWaiting {
		return nil
	}
	displaySvc.PrintInfo("⏳ Waiting for opponent to join (press 'q' to leave)...")

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case line, ok := <-input:
			if !ok {
				input = nil // stdin closed; keep polling for an opponent
				continue
			}
			if _, _, giveUp, _ := lib.ParseMove(line); giveUp {
				if _, err := gameSvc.Quit(ctx, lib.GameState.GetID(), lib.UserID); err != nil {
					return fmt.Errorf("failed to quit: %w", err)
				}
				return ErrUserQuit
			}

		case <-ticker.C:
			g, err := gameSvc.GetGame(ctx, lib.GameState.GetID())
			if err != nil {
				continue // transient; retry on the next tick
			}
			if g.GetStatus() == openapi.StatusInProgress {
				lib.GameState = g
				return nil
			}
		}
	}
}
