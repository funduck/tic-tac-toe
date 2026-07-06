package game

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"
)

// TestConcurrentJoinGame verifies that when multiple players race to fill the
// last open slot in a waiting game, exactly one succeeds and the game reaches
// StatusInProgress with a single UserID2.
func TestConcurrentJoinGame(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepo()
	svc := NewGameService(repo, slog.Default())

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
	svc := NewGameService(repo, slog.Default())

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

// TestConcurrentPollAndMove is a lost-update regression test: GetGame persists
// presence (read-modify-write), so a poll racing a MakeMove must never write
// stale state back and silently erase the move.
func TestConcurrentPollAndMove(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepo()
	svc := NewGameService(repo, slog.Default())

	g, _ := svc.CreateGame(ctx, "alice", CreateGameCommand{})
	_, _ = svc.JoinGame(ctx, "alice", JoinGameCommand{GameID: g.ID})
	_, _ = svc.JoinGame(ctx, "bob", JoinGameCommand{GameID: g.ID})

	// Alice takes the first column, bob fills the second; alice wins.
	moves := []struct {
		player string
		x, y   int
	}{
		{"alice", 0, 0}, {"bob", 1, 0},
		{"alice", 0, 1}, {"bob", 1, 1},
		{"alice", 0, 2},
	}

	for _, mv := range moves {
		opponent := "alice"
		if mv.player == "alice" {
			opponent = "bob"
		}

		// The opponent polls aggressively while the move is being made.
		done := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					if _, err := svc.GetGame(ctx, opponent, GetGameCommand{GameID: g.ID}); err != nil {
						t.Errorf("unexpected poll error: %v", err)
						return
					}
				}
			}
		}()

		if _, err := svc.MakeMove(ctx, mv.player, MakeMoveCommand{GameID: g.ID, X: mv.x, Y: mv.y}); err != nil {
			t.Fatalf("move (%d,%d) by %s failed: %v", mv.x, mv.y, mv.player, err)
		}
		close(done)
		wg.Wait()

		final, err := svc.GetGame(ctx, mv.player, GetGameCommand{GameID: g.ID})
		if err != nil {
			t.Fatalf("failed to get game: %v", err)
		}
		if final.Board[mv.x][mv.y] == 0 {
			t.Fatalf("move (%d,%d) by %s was lost to a concurrent poll", mv.x, mv.y, mv.player)
		}
	}

	final, err := svc.GetGame(ctx, "alice", GetGameCommand{GameID: g.ID})
	if err != nil {
		t.Fatalf("failed to get game: %v", err)
	}
	if final.Status != StatusFinished || final.WinnerID != "alice" {
		t.Errorf("expected alice to win, got status=%s winner=%q", final.Status, final.WinnerID)
	}
}
