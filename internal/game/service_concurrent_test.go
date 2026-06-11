package game

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
)

// TestConcurrentGames validates that multiple games can be played simultaneously
// without race conditions or data corruption. This test creates 3 games with
// 6 different players (2 per game) and plays them concurrently.
func TestConcurrentGames(t *testing.T) {
	// Setup: create service with shared repository
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Reduce noise during test
	}))
	repo := NewMemoryRepo()
	service := NewGameService(repo, logger)

	// Number of concurrent games to play
	const numGames = 3

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Channel to collect errors from goroutines
	errChan := make(chan error, numGames)

	// Channel to collect results
	results := make(chan *Game, numGames)

	// Start numGames concurrent games
	for i := 0; i < numGames; i++ {
		wg.Add(1)
		gameNum := i + 1

		go func(gameIndex int) {
			defer wg.Done()

			// Unique players for this game
			player1 := fmt.Sprintf("player_%d_1", gameIndex)
			player2 := fmt.Sprintf("player_%d_2", gameIndex)

			// Create game
			game, err := service.CreateGame()
			if err != nil {
				errChan <- fmt.Errorf("game %d: failed to create game: %w", gameIndex, err)
				return
			}

			// Player 1 joins
			game, err = service.JoinGame(game.ID, player1)
			if err != nil {
				errChan <- fmt.Errorf("game %d: player 1 failed to join: %w", gameIndex, err)
				return
			}

			// Player 2 joins
			game, err = service.JoinGame(game.ID, player2)
			if err != nil {
				errChan <- fmt.Errorf("game %d: player 2 failed to join: %w", gameIndex, err)
				return
			}

			// Verify game is in progress
			if game.Status != StatusInProgress {
				errChan <- fmt.Errorf("game %d: expected status %s, got %s", gameIndex, StatusInProgress, game.Status)
				return
			}

			// Play a simple game sequence
			// Player 1 (X) wins with top row: X X X
			moves := []struct {
				player string
				x, y   int
			}{
				{player1, 0, 0}, // X at (0,0)
				{player2, 1, 0}, // O at (1,0)
				{player1, 0, 1}, // X at (0,1)
				{player2, 1, 1}, // O at (1,1)
				{player1, 0, 2}, // X at (0,2) - wins!
			}

			for _, move := range moves {
				game, err = service.MakeMove(game.ID, move.player, move.x, move.y)
				if err != nil {
					errChan <- fmt.Errorf("game %d: move failed at (%d,%d): %w", gameIndex, move.x, move.y, err)
					return
				}
			}

			// Verify game finished with player 1 as winner
			if game.Status != StatusFinished {
				errChan <- fmt.Errorf("game %d: expected status %s, got %s", gameIndex, StatusFinished, game.Status)
				return
			}

			if game.Result != ResultWin {
				errChan <- fmt.Errorf("game %d: expected result %s, got %s", gameIndex, ResultWin, game.Result)
				return
			}

			if game.WinnerID != player1 {
				errChan <- fmt.Errorf("game %d: expected winner %s, got %s", gameIndex, player1, game.WinnerID)
				return
			}

			// Send result
			results <- game

		}(gameNum)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)
	close(results)

	// Check for errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		t.Errorf("Concurrent game execution failed with %d error(s):", len(errors))
		for _, err := range errors {
			t.Errorf("  - %v", err)
		}
		return
	}

	// Verify we got all results
	completedGames := 0
	for range results {
		completedGames++
	}

	if completedGames != numGames {
		t.Errorf("Expected %d completed games, got %d", numGames, completedGames)
	}

	t.Logf("Successfully completed %d concurrent games", completedGames)
}

// TestConcurrentJoinAnyGame tests the matchmaking functionality under concurrent load.
// It verifies that multiple players can simultaneously join waiting games without conflicts.
func TestConcurrentJoinAnyGame(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	repo := NewMemoryRepo()
	service := NewGameService(repo, logger)

	// Create 2 waiting games
	game1, err := service.CreateGame()
	if err != nil {
		t.Fatalf("Failed to create game 1: %v", err)
	}

	game2, err := service.CreateGame()
	if err != nil {
		t.Fatalf("Failed to create game 2: %v", err)
	}

	// Player 1 joins game 1
	_, err = service.JoinGame(game1.ID, "player1")
	if err != nil {
		t.Fatalf("Player 1 failed to join game 1: %v", err)
	}

	// Player 2 joins game 2
	_, err = service.JoinGame(game2.ID, "player2")
	if err != nil {
		t.Fatalf("Player 2 failed to join game 2: %v", err)
	}

	// Now we have 2 waiting games, each with 1 player
	// Try to join both concurrently with 2 new players
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	gamesChan := make(chan *Game, 2)

	// Player 3 joins any game
	wg.Add(1)
	go func() {
		defer wg.Done()
		game, err := service.JoinAnyGame("player3")
		if err != nil {
			errChan <- fmt.Errorf("player3 failed to join: %w", err)
			return
		}
		gamesChan <- game
	}()

	// Player 4 joins any game
	wg.Add(1)
	go func() {
		defer wg.Done()
		game, err := service.JoinAnyGame("player4")
		if err != nil {
			errChan <- fmt.Errorf("player4 failed to join: %w", err)
			return
		}
		gamesChan <- game
	}()

	wg.Wait()
	close(errChan)
	close(gamesChan)

	// Check for errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		t.Errorf("Concurrent join failed with errors:")
		for _, err := range errors {
			t.Errorf("  - %v", err)
		}
		return
	}

	// Verify both games are now in progress
	joinedGames := make(map[string]bool)
	for game := range gamesChan {
		if game.Status != StatusInProgress {
			t.Errorf("Expected game %s to be in progress, got status: %s", game.ID, game.Status)
		}
		joinedGames[game.ID] = true
	}

	if len(joinedGames) != 2 {
		t.Errorf("Expected 2 different games to be joined, got %d", len(joinedGames))
	}

	t.Logf("Successfully matched %d players to %d games concurrently", 2, len(joinedGames))
}

// TestRepositoryConcurrentAccess validates thread-safety of the repository layer
// by performing many concurrent reads and writes.
func TestRepositoryConcurrentAccess(t *testing.T) {
	repo := NewMemoryRepo()

	// Create initial game
	game := NewGame("test-game-123")
	game.UserID1 = "player1"
	game.UserID2 = "player2"
	game.Status = StatusInProgress
	game.CurrentPlayerID = "player1"

	err := repo.Create(game)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	const numReaders = 10
	const numWriters = 5
	const operationsPerGoroutine = 20

	var wg sync.WaitGroup
	errChan := make(chan error, numReaders+numWriters)

	// Start concurrent readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				_, err := repo.FindByID("test-game-123")
				if err != nil {
					errChan <- fmt.Errorf("reader %d: read failed: %w", readerID, err)
					return
				}
			}
		}(i)
	}

	// Start concurrent writers (making moves)
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				g, err := repo.FindByID("test-game-123")
				if err != nil {
					errChan <- fmt.Errorf("writer %d: read failed: %w", writerID, err)
					return
				}

				// Toggle current player (simulating moves)
				if g.CurrentPlayerID == "player1" {
					g.CurrentPlayerID = "player2"
				} else {
					g.CurrentPlayerID = "player1"
				}

				err = repo.Update(g)
				if err != nil {
					errChan <- fmt.Errorf("writer %d: update failed: %w", writerID, err)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		t.Errorf("Concurrent repository access failed with %d error(s):", len(errors))
		for _, err := range errors {
			t.Errorf("  - %v", err)
		}
		return
	}

	// Verify final state is consistent
	finalGame, err := repo.FindByID("test-game-123")
	if err != nil {
		t.Fatalf("Failed to read final game state: %v", err)
	}

	if finalGame.ID != "test-game-123" {
		t.Errorf("Game ID corrupted: expected 'test-game-123', got '%s'", finalGame.ID)
	}

	if finalGame.CurrentPlayerID != "player1" && finalGame.CurrentPlayerID != "player2" {
		t.Errorf("Current player ID corrupted: got '%s'", finalGame.CurrentPlayerID)
	}

	totalOps := numReaders*operationsPerGoroutine + numWriters*operationsPerGoroutine
	t.Logf("Successfully completed %d concurrent repository operations (%d reads, %d writes)",
		totalOps, numReaders*operationsPerGoroutine, numWriters*operationsPerGoroutine)
}
