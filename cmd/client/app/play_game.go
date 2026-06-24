package app

import (
	"context"
	"fmt"
	"time"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/funduck/tic-tac-toe/client/lib"
)

// pollInterval is how often the client refreshes game state while waiting.
const pollInterval = time.Second

// PlayGame runs the main game loop. It polls for opponent moves on a ticker
// while concurrently reading stdin, so the player can make a move on their turn
// or forfeit with 'q' at any time.
func PlayGame(ctx context.Context, gameSvc *lib.GameService, displaySvc *lib.DisplayService, input <-chan string) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	render(displaySvc)

	for lib.GameState.GetStatus() != openapi.StatusFinished {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case line, ok := <-input:
			if !ok {
				input = nil // stdin closed; keep polling for the result
				continue
			}
			done, err := handleInput(ctx, gameSvc, displaySvc, line)
			if err != nil {
				return err
			}
			if done {
				return nil
			}

		case <-ticker.C:
			g, err := gameSvc.GetGame(ctx, lib.GameState.GetID())
			if err != nil {
				if lib.IsSessionError(err) {
					return err
				}
				continue // transient; retry on the next tick
			}
			applyUpdate(g, displaySvc)
		}
	}

	return nil
}

// handleInput processes one line of player input. It returns done=true when the
// game loop should exit (the player forfeited).
func handleInput(ctx context.Context, gameSvc *lib.GameService, displaySvc *lib.DisplayService, line string) (done bool, err error) {
	row, col, giveUp, parseErr := lib.ParseMove(line)

	if giveUp {
		g, err := gameSvc.GiveUp(ctx, lib.GameState.GetID(), lib.UserID)
		if err != nil {
			return false, fmt.Errorf("failed to give up: %w", err)
		}
		applyUpdate(g, displaySvc)
		return true, nil
	}

	if lib.GameState.GetCurrentPlayerID() != lib.UserID {
		displaySvc.PrintError("Please wait for your turn (or 'q' to quit)")
		return false, nil
	}

	if parseErr != nil {
		displaySvc.PrintError(parseErr.Error())
		return false, nil
	}

	g, err := gameSvc.MakeMove(ctx, lib.GameState.GetID(), lib.UserID, row, col)
	if err != nil {
		if lib.IsSessionError(err) {
			return false, err
		}
		displaySvc.PrintError(lib.FriendlyMessage(err))
		// Refresh state in case the board moved on while we were typing.
		if fresh, gerr := gameSvc.GetGame(ctx, lib.GameState.GetID()); gerr == nil {
			applyUpdate(fresh, displaySvc)
		}
		return false, nil
	}
	applyUpdate(g, displaySvc)
	return false, nil
}

// applyUpdate swaps in new game state and redraws only if something changed,
// avoiding needless flicker on no-op polls.
func applyUpdate(g *lib.Game, displaySvc *lib.DisplayService) {
	if sameGame(lib.GameState, g) {
		return
	}
	lib.LastMove = diffMove(lib.GameState.GetBoard(), g.GetBoard())
	lib.GameState = g
	render(displaySvc)
}

// render draws the current frame and, while the game is in progress, a prompt
// reminding the player what they can type.
func render(displaySvc *lib.DisplayService) {
	displaySvc.RenderFrame()
	if lib.GameState.GetStatus() != openapi.StatusInProgress {
		return
	}
	if lib.GameState.GetCurrentPlayerID() == lib.UserID {
		fmt.Print("Your move (1-9, or 'q' to quit): ")
	} else {
		fmt.Print("Waiting for opponent (type 'q' to quit)... ")
	}
}

// sameGame reports whether two states are visually identical for our purposes.
func sameGame(a, b *lib.Game) bool {
	if a.GetStatus() != b.GetStatus() || a.GetCurrentPlayerID() != b.GetCurrentPlayerID() {
		return false
	}
	ab, bb := a.GetBoard(), b.GetBoard()
	for r := range 3 {
		for c := range 3 {
			if cellAt(ab, r, c) != cellAt(bb, r, c) {
				return false
			}
		}
	}
	return true
}

// diffMove returns the 1-based numpad cell that newly gained a mark between the
// two boards, or -1 if none.
func diffMove(prev, next [][]int32) int {
	for r := range 3 {
		for c := range 3 {
			if cellAt(prev, r, c) == 0 && cellAt(next, r, c) != 0 {
				return r*3 + c + 1
			}
		}
	}
	return -1
}

func cellAt(board [][]int32, r, c int) int32 {
	if r < len(board) && c < len(board[r]) {
		return board[r][c]
	}
	return 0
}
