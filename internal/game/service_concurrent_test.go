package game

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"
)

// TestConcurrentJoinGame verifies that when multiple players race to fill the
// last open slot in a waiting game, exactly one succeeds and the game reaches
// StatusInProgress with a single UserID2.
func TestConcurrentJoinGame(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)

	g, _ := svc.CreateGame(ctx, "host", CreateGameCommand{})
	_, _ = svc.JoinGame(ctx, "alice", JoinGameCommand{GameID: g.ID})

	const racers = 5
	type result struct {
		game *Game
		err  error
	}
	results := make(chan result, racers)

	var wg sync.WaitGroup
	for i := range racers {
		wg.Add(1)
		playerID := fmt.Sprintf("racer_%d", i)
		go func(id string) {
			defer wg.Done()
			g, err := svc.JoinGame(ctx, id, JoinGameCommand{GameID: g.ID})
			results <- result{g, err}
		}(playerID)
	}
	wg.Wait()
	close(results)

	var succeeded int
	for r := range results {
		if r.err == nil {
			succeeded++
		}
	}

	if succeeded != 1 {
		t.Errorf("expected exactly 1 successful join, got %d", succeeded)
	}

	final, err := svc.GetGame(ctx, "alice", GetGameCommand{GameID: g.ID})
	if err != nil {
		t.Fatalf("failed to get game: %v", err)
	}
	if final.Status != StatusInProgress {
		t.Errorf("expected status in_progress, got %s", final.Status)
	}
	if final.UserID2 == "" {
		t.Error("expected UserID2 to be set")
	}
}

// TestConcurrentJoinAnyGame verifies that concurrent players joining via
// matchmaking each land in a different waiting game without conflicts.
func TestConcurrentJoinAnyGame(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)

	game1, _ := svc.CreateGame(ctx, "host1", CreateGameCommand{})
	game2, _ := svc.CreateGame(ctx, "host2", CreateGameCommand{})
	_, _ = svc.JoinGame(ctx, "player1", JoinGameCommand{GameID: game1.ID})
	_, _ = svc.JoinGame(ctx, "player2", JoinGameCommand{GameID: game2.ID})

	// player3 and player4 race to fill the two open slots
	type result struct {
		game *Game
		err  error
	}
	results := make(chan result, 2)

	var wg sync.WaitGroup
	for _, id := range []string{"player3", "player4"} {
		wg.Add(1)
		go func(playerID string) {
			defer wg.Done()
			g, err := svc.JoinAnyGame(ctx, playerID)
			results <- result{g, err}
		}(id)
	}
	wg.Wait()
	close(results)

	joinedGames := make(map[string]bool)
	for r := range results {
		if r.err != nil {
			t.Errorf("unexpected error: %v", r.err)
			continue
		}
		if r.game.Status != StatusInProgress {
			t.Errorf("expected game %s in_progress, got %s", r.game.ID, r.game.Status)
		}
		joinedGames[r.game.ID] = true
	}

	if len(joinedGames) != 2 {
		t.Errorf("expected players in 2 different games, got %d", len(joinedGames))
	}
}
