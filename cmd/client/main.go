package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/funduck/tic-tac-toe/client/app"
	"github.com/funduck/tic-tac-toe/client/lib"
	"golang.org/x/term"
)

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "server base URL")
	userID := flag.String("user", "", "your user ID (required)")
	password := flag.String("password", "", "your password (prompted if omitted)")
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
		pw, err := promptPassword()
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌ Error:", err)
			os.Exit(1)
		}
		lib.Password = pw
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

// promptPassword reads a password from the terminal without echoing it. It
// fails when stdin is not interactive so scripted usage still gets a clear error.
func promptPassword() (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return "", fmt.Errorf("--password is required when stdin is not a terminal")
	}
	fmt.Print("Password: ")
	pw, err := term.ReadPassword(fd)
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	if len(pw) == 0 {
		return "", fmt.Errorf("password cannot be empty")
	}
	return string(pw), nil
}
