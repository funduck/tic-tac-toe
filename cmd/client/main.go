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
	password := flag.String("password", "", "your password (required)")
	gameID := flag.String("game", "", "game ID to join (omit to create a new game)")
	private := flag.Bool("private", false, "create a private game")

	flag.Parse()

	lib.UserID = *userID
	lib.Password = *password
	lib.GameID = *gameID
	lib.Private = *private

	if lib.UserID == "" {
		fmt.Fprintln(os.Stderr, "❌ Error: --user is required")
		flag.Usage()
		os.Exit(1)
	}
	if lib.Password == "" {
		fmt.Fprintln(os.Stderr, "❌ Error: --password is required")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize API client configuration
	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: *serverURL}}

	if os.Getenv("DEBUG") != "" {
		cfg.Debug = true
	}

	// Create a temporary API client for UserService initialization
	tempAPIClient := openapi.NewAPIClient(cfg)
	userSvc := lib.NewUserService(tempAPIClient)

	// Create authenticated HTTP client
	authHTTPClient := lib.NewAuthHTTPClient(userSvc)
	cfg.HTTPClient = authHTTPClient

	// Create final API client with authentication
	apiClient := openapi.NewAPIClient(cfg)

	scanner := bufio.NewScanner(os.Stdin)

	// Initialize services
	gameSvc := lib.NewGameService(apiClient)
	displaySvc := lib.NewDisplayService()
	inputSvc := lib.NewInputService(scanner)

	ctx := context.Background()

	app.Start(ctx, gameSvc, userSvc, displaySvc, inputSvc)
}
