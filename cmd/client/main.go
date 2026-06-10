package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"

	"github.com/funduck/tic-tac-toe/cmd/client/app"
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

	scanner := bufio.NewScanner(os.Stdin)

	// Initialize services
	gameSvc := lib.NewGameService(apiClient)
	displaySvc := lib.NewDisplayService()
	inputSvc := lib.NewInputService(scanner)

	ctx := context.Background()

	app.Start(ctx, gameSvc, displaySvc, inputSvc)
}
