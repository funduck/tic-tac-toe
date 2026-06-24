package game

import (
	"errors"
	"testing"
	"time"
)

var (
	testAfkTimeout = 30 * time.Second
)

func TestNewGame(t *testing.T) {
	g := NewGame("id1")
	if g.ID != "id1" {
		t.Errorf("expected ID id1, got %s", g.ID)
	}
	if g.Status != StatusWaiting {
		t.Errorf("expected status %s, got %s", StatusWaiting, g.Status)
	}
	if g.Board != [3][3]int{} {
		t.Error("expected empty board")
	}
}

func TestNewGameInProgress(t *testing.T) {
	g := NewGameInProgress("id1", "alice", "bob")
	if g.Status != StatusInProgress {
		t.Errorf("expected status %s, got %s", StatusInProgress, g.Status)
	}
	if g.CurrentPlayerID != "alice" {
		t.Errorf("expected currentPlayer alice, got %s", g.CurrentPlayerID)
	}
	if g.Board != [3][3]int{} {
		t.Error("expected empty board")
	}
}

func TestForfeitIfOpponentAFK(t *testing.T) {
	now := time.Now()

	t.Run("opponent stale, active player wins", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		g.Touch("alice", now)
		g.Touch("bob", now.Add(-testAfkTimeout-time.Second))

		if !g.ForfeitIfOpponentAFK("alice", now, testAfkTimeout) {
			t.Fatal("expected forfeit to fire")
		}
		if g.Status != StatusFinished || g.Result != ResultWin || g.WinnerID != "alice" {
			t.Errorf("unexpected end state: status=%s result=%s winner=%s", g.Status, g.Result, g.WinnerID)
		}
	})

	t.Run("opponent recently seen, no change", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		g.Touch("alice", now)
		g.Touch("bob", now.Add(-time.Second))

		if g.ForfeitIfOpponentAFK("alice", now, testAfkTimeout) {
			t.Fatal("did not expect forfeit")
		}
		if g.Status != StatusInProgress {
			t.Errorf("expected status still in progress, got %s", g.Status)
		}
	})

	t.Run("opponent never seen, no change", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		g.Touch("alice", now)

		if g.ForfeitIfOpponentAFK("alice", now, testAfkTimeout) {
			t.Fatal("did not expect forfeit when opponent never seen")
		}
	})

	t.Run("not in progress, no change", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		g.Status = StatusFinished
		g.Touch("bob", now.Add(-testAfkTimeout-time.Second))

		if g.ForfeitIfOpponentAFK("alice", now, testAfkTimeout) {
			t.Fatal("did not expect forfeit on a non-in-progress game")
		}
	})

	t.Run("non-participant, no change", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		g.Touch("bob", now.Add(-testAfkTimeout-time.Second))

		if g.ForfeitIfOpponentAFK("carol", now, testAfkTimeout) {
			t.Fatal("did not expect forfeit for a non-participant")
		}
	})
}

func TestJoin_Valid(t *testing.T) {
	t.Run("first player joins", func(t *testing.T) {
		g := NewGame("id1")
		if err := g.Join("alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if g.UserID1 != "alice" {
			t.Errorf("expected UserID1 alice, got %s", g.UserID1)
		}
		if g.Status != StatusWaiting {
			t.Errorf("expected status still waiting, got %s", g.Status)
		}
	})

	t.Run("second player joins, game starts", func(t *testing.T) {
		g := NewGame("id1")
		if err := g.Join("alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := g.Join("bob"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if g.UserID2 != "bob" {
			t.Errorf("expected UserID2 bob, got %s", g.UserID2)
		}
		if g.Status != StatusInProgress {
			t.Errorf("expected status in_progress, got %s", g.Status)
		}
	})

	t.Run("same player joins again", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		// allow join again
		if err := g.Join("alice"); err != nil {
			t.Fatalf("unexpected error on rejoin: %v", err)
		}
		if g.UserID1 != "alice" {
			t.Errorf("expected UserID1 alice, got %s", g.UserID1)
		}
		if g.UserID2 != "bob" {
			t.Errorf("expected UserID2 bob, got %s", g.UserID2)
		}
		if g.Status != StatusInProgress {
			t.Errorf("expected status still in_progress, got %s", g.Status)
		}
	})
}

func TestJoin_Errors(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Game)
		userID  string
		wantErr error
	}{
		{
			name:    "join after game started",
			setup:   func(g *Game) { g.UserID1 = "alice"; g.UserID2 = "bob"; g.Status = StatusInProgress },
			userID:  "charlie",
			wantErr: ErrGameNotWaiting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame("id1")
			tt.setup(g)
			err := g.Join(tt.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestMakeMove_Valid(t *testing.T) {
	g := NewGameInProgress("id1", "alice", "bob")

	if err := g.MakeMove("alice", 0, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.Board[0][0] != 1 {
		t.Errorf("expected board[0][0]=1, got %d", g.Board[0][0])
	}
	if g.CurrentPlayerID != "bob" {
		t.Errorf("expected next player bob, got %s", g.CurrentPlayerID)
	}
	if g.Status != StatusInProgress {
		t.Errorf("expected game still in progress, got %s", g.Status)
	}

	if err := g.MakeMove("bob", 1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.Board[1][1] != 2 {
		t.Errorf("expected board[1][1]=2, got %d", g.Board[1][1])
	}
	if g.CurrentPlayerID != "alice" {
		t.Errorf("expected next player alice, got %s", g.CurrentPlayerID)
	}
}

func TestMakeMove_Errors(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Game)
		userID  string
		x, y    int
		wantErr error
	}{
		{
			name:   "out of bounds negative x",
			setup:  func(g *Game) {},
			userID: "alice", x: -1, y: 0,
			wantErr: ErrOutOfBounds,
		},
		{
			name:   "out of bounds x > 2",
			setup:  func(g *Game) {},
			userID: "alice", x: 3, y: 0,
			wantErr: ErrOutOfBounds,
		},
		{
			name:   "out of bounds y > 2",
			setup:  func(g *Game) {},
			userID: "alice", x: 0, y: 3,
			wantErr: ErrOutOfBounds,
		},
		{
			name:   "cell occupied",
			setup:  func(g *Game) { g.Board[1][1] = 1 },
			userID: "alice", x: 1, y: 1,
			wantErr: ErrCellOccupied,
		},
		{
			name:   "not your turn",
			setup:  func(g *Game) {},
			userID: "bob", x: 0, y: 0,
			wantErr: ErrNotYourTurn,
		},
		{
			name:   "game not active",
			setup:  func(g *Game) { g.Status = StatusFinished },
			userID: "alice", x: 0, y: 0,
			wantErr: ErrGameNotActive,
		},
		{
			name:   "user not in game",
			setup:  func(g *Game) {},
			userID: "charlie", x: 0, y: 0,
			wantErr: ErrNotInGame,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGameInProgress("id1", "alice", "bob")
			tt.setup(g)
			err := g.MakeMove(tt.userID, tt.x, tt.y)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// TestMakeMove_WinConditions covers all 8 winning lines.
func TestMakeMove_WinConditions(t *testing.T) {
	winCases := []struct {
		name       string
		aliceMoves [][2]int
	}{
		{"row0", [][2]int{{0, 0}, {0, 1}, {0, 2}}},
		{"row1", [][2]int{{1, 0}, {1, 1}, {1, 2}}},
		{"row2", [][2]int{{2, 0}, {2, 1}, {2, 2}}},
		{"col0", [][2]int{{0, 0}, {1, 0}, {2, 0}}},
		{"col1", [][2]int{{0, 1}, {1, 1}, {2, 1}}},
		{"col2", [][2]int{{0, 2}, {1, 2}, {2, 2}}},
		{"diag", [][2]int{{0, 0}, {1, 1}, {2, 2}}},
		{"antidiag", [][2]int{{0, 2}, {1, 1}, {2, 0}}},
	}

	for _, wc := range winCases {
		t.Run(wc.name, func(t *testing.T) {
			g := NewGameInProgress("id1", "alice", "bob")

			// collect alice's cells
			aliceSet := map[[2]int]bool{}
			for _, m := range wc.aliceMoves {
				aliceSet[[2]int{m[0], m[1]}] = true
			}

			// pick bob filler cells (won't interfere with alice)
			var bobMoves [][2]int
			for r := 0; r < 3 && len(bobMoves) < len(wc.aliceMoves)-1; r++ {
				for c := 0; c < 3 && len(bobMoves) < len(wc.aliceMoves)-1; c++ {
					if !aliceSet[[2]int{r, c}] {
						bobMoves = append(bobMoves, [2]int{r, c})
					}
				}
			}

			// interleave: alice, bob, alice, bob, alice
			for i, m := range wc.aliceMoves {
				if err := g.MakeMove("alice", m[0], m[1]); err != nil {
					t.Fatalf("alice move %v failed: %v", m, err)
				}
				if g.Status == StatusFinished {
					break
				}
				if i < len(bobMoves) {
					if err := g.MakeMove("bob", bobMoves[i][0], bobMoves[i][1]); err != nil {
						t.Fatalf("bob move %v failed: %v", bobMoves[i], err)
					}
				}
			}

			if g.Status != StatusFinished {
				t.Error("expected game to be finished")
			}
			if g.Result != ResultWin {
				t.Errorf("expected result win, got %s", g.Result)
			}
			if g.WinnerID != "alice" {
				t.Errorf("expected winner alice, got %s", g.WinnerID)
			}
		})
	}
}

// TestMakeMove_Draw uses a known draw sequence:
//
//	X O X
//	X X O
//	O X O
func TestMakeMove_Draw(t *testing.T) {
	g := NewGameInProgress("id1", "alice", "bob")
	moves := []struct {
		user string
		x, y int
	}{
		{"alice", 0, 0}, // X
		{"bob", 0, 1},   // O
		{"alice", 0, 2}, // X
		{"bob", 1, 2},   // O
		{"alice", 1, 0}, // X
		{"bob", 2, 0},   // O
		{"alice", 1, 1}, // X
		{"bob", 2, 2},   // O
		{"alice", 2, 1}, // X  → board full, no winner
	}
	for _, m := range moves {
		if err := g.MakeMove(m.user, m.x, m.y); err != nil {
			t.Fatalf("move (%s,%d,%d) failed: %v", m.user, m.x, m.y, err)
		}
	}
	if g.Status != StatusFinished {
		t.Error("expected game to be finished")
	}
	if g.Result != ResultDraw {
		t.Errorf("expected draw, got %s", g.Result)
	}
	if g.WinnerID != "" {
		t.Errorf("expected no winner, got %s", g.WinnerID)
	}
}

func TestGiveUp(t *testing.T) {
	t.Run("player1 gives up, player2 wins", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		if err := g.GiveUp("alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if g.Status != StatusFinished {
			t.Errorf("expected finished, got %s", g.Status)
		}
		if g.Result != ResultWin {
			t.Errorf("expected win result, got %s", g.Result)
		}
		if g.WinnerID != "bob" {
			t.Errorf("expected winner bob, got %s", g.WinnerID)
		}
	})

	t.Run("player2 gives up, player1 wins", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		if err := g.GiveUp("bob"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if g.WinnerID != "alice" {
			t.Errorf("expected winner alice, got %s", g.WinnerID)
		}
	})

	t.Run("give up on finished game", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		g.Status = StatusFinished
		if err := g.GiveUp("alice"); !errors.Is(err, ErrGameNotActive) {
			t.Errorf("expected ErrGameNotActive, got %v", err)
		}
	})

	t.Run("non-player gives up", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		if err := g.GiveUp("charlie"); !errors.Is(err, ErrNotInGame) {
			t.Errorf("expected ErrNotInGame, got %v", err)
		}
	})
}

func TestQuit(t *testing.T) {
	t.Run("creator quits waiting game, it is cancelled", func(t *testing.T) {
		g := NewGame("id1")
		if err := g.Join("alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := g.Quit("alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if g.Status != StatusCancelled {
			t.Errorf("expected cancelled, got %s", g.Status)
		}
	})

	t.Run("quit on in-progress game is rejected", func(t *testing.T) {
		g := NewGameInProgress("id1", "alice", "bob")
		if err := g.Quit("alice"); !errors.Is(err, ErrGameNotWaiting) {
			t.Errorf("expected ErrGameNotWaiting, got %v", err)
		}
	})

	t.Run("non-player quits", func(t *testing.T) {
		g := NewGame("id1")
		if err := g.Join("alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := g.Quit("charlie"); !errors.Is(err, ErrNotInGame) {
			t.Errorf("expected ErrNotInGame, got %v", err)
		}
	})
}

func TestGame_Touch(t *testing.T) {
	now := time.Now()
	g := NewGameInProgress("id1", "alice", "bob")

	t.Run("stamps first player", func(t *testing.T) {
		if !g.Touch("alice", now) {
			t.Fatal("expected Touch to report a change for a participant")
		}
		if g.UserID1LastSeen == nil || !g.UserID1LastSeen.Equal(now) {
			t.Error("expected UserID1LastSeen to be set to now")
		}
		if g.UserID2LastSeen != nil {
			t.Error("touching alice should not stamp bob")
		}
	})

	t.Run("stamps second player", func(t *testing.T) {
		if !g.Touch("bob", now) {
			t.Fatal("expected Touch to report a change for a participant")
		}
		if g.UserID2LastSeen == nil || !g.UserID2LastSeen.Equal(now) {
			t.Error("expected UserID2LastSeen to be set to now")
		}
	})

	t.Run("ignores non-participant", func(t *testing.T) {
		if g.Touch("charlie", now) {
			t.Error("expected Touch to report no change for a non-participant")
		}
	})

	t.Run("ignores empty user", func(t *testing.T) {
		empty := NewGame("id2") // UserID1 == "" but empty userID must not match
		if empty.Touch("", now) {
			t.Error("expected Touch to ignore an empty userID")
		}
	})
}
