package app

import (
	"context"
	"fmt"
	"os"

	"github.com/funduck/tic-tac-toe/cmd/client/lib"
)

func Start(ctx context.Context, gameSvc *lib.GameService, displaySvc *lib.DisplayService, inputSvc *lib.InputService) {
	// Create or join game
	err := CreateOrJoinGame(ctx, gameSvc, displaySvc)
	if err != nil {
		os.Exit(1)
	}

	displaySvc.PrintInfo(fmt.Sprintf("Game started! You are %s.", displaySvc.MyMark()))
	displaySvc.PrintBoard()

	// Game loop
	err = PlayGame(ctx, gameSvc, displaySvc, inputSvc)
	if err != nil {
		os.Exit(1)
	}

	// Display result
	displaySvc.PrintResult()
}
