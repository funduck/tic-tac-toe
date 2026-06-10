package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"

	"github.com/funduck/tic-tac-toe/cmd/client/lib"
)

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "server base URL")
	userID := flag.String("user", "", "your user ID (required)")
	gameID := flag.String("game", "", "game ID to join (omit to create a new game)")
	flag.Parse()
	lib.UserID = *userID
	lib.GameID = *gameID

	if lib.UserID == "" {
		fmt.Fprintln(os.Stderr, "❌ Error: --user is required")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize API client
	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: *serverURL}}
	apiClient := openapi.NewAPIClient(cfg)

	// Initialize services
	gameSvc := lib.NewGameService(apiClient)
	displaySvc := lib.NewDisplayService()
	inputSvc := lib.NewInputService()

	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	// Create or join game
	err := lib.CreateOrJoinGame(ctx, gameSvc, displaySvc)
	if err != nil {
		displaySvc.PrintError(fmt.Sprintf("Failed to start game: %v", err))
		os.Exit(1)
	}

	displaySvc.PrintInfo(fmt.Sprintf("Game started! You are %s.", displaySvc.MyMark()))
	displaySvc.PrintBoard()

	// Game loop
	err = lib.PlayGame(ctx, gameSvc, displaySvc, inputSvc, scanner)
	if err != nil {
		displaySvc.PrintError(fmt.Sprintf("Game error: %v", err))
		os.Exit(1)
	}

	// Display result
	displaySvc.PrintResult()
}
