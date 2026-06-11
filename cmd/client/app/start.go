package app

import (
	"context"
	"fmt"
	"os"

	"github.com/funduck/tic-tac-toe/client/lib"
)

func Start(ctx context.Context, gameSvc *lib.GameService, userSvc *lib.UserService, displaySvc *lib.DisplayService, inputSvc *lib.InputService) {
	// Authenticate user (login or signup)
	err := LoginOrSignup(ctx, userSvc, displaySvc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Authentication failed:", err)
		os.Exit(1)
	}

	// Create or join game
	err = CreateOrJoinGame(ctx, gameSvc, displaySvc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error:", err)
		os.Exit(1)
	}

	// Wait for opponent if needed
	err = WaitForOpponent(ctx, gameSvc, displaySvc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error:", err)
		os.Exit(1)
	}

	displaySvc.PrintInfo(fmt.Sprintf("Game started! You are %s.", displaySvc.MyMark()))
	displaySvc.PrintBoard()

	// Game loop
	err = PlayGame(ctx, gameSvc, displaySvc, inputSvc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error:", err)
		os.Exit(1)
	}

	// Display result
	displaySvc.PrintResult()
}
