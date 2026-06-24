package game

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

var ErrMockRepoUpdateFailure = errors.New("update failure")

// mockRepo is a hand-written test double for GameRepo.
type mockRepo struct {
	games     map[string]*Game
	order     []string
	createErr error
	getErr    error
	updateErr error
}

func newMockRepo() *mockRepo {
	return &mockRepo{games: make(map[string]*Game)}
}

func (m *mockRepo) Create(ctx context.Context, game *Game) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.games[game.ID] = game
	m.order = append(m.order, game.ID)
	return nil
}

func (m *mockRepo) FindByID(ctx context.Context, gameID string) (*Game, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	g, ok := m.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}
	return g, nil
}

func (m *mockRepo) Update(ctx context.Context, game *Game) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.games[game.ID]; !ok {
		return ErrGameNotFound
	}
	m.games[game.ID] = game
	return nil
}

func (m *mockRepo) FindLatestForUser(ctx context.Context, userID string) (*Game, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for i := len(m.order) - 1; i >= 0; i-- {
		g := m.games[m.order[i]]
		if g.UserID1 == userID || g.UserID2 == userID {
			return g, nil
		}
	}
	return nil, ErrGameNotFound
}

func (m *mockRepo) FindGameToJoin(ctx context.Context, userID string) (*Game, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for i := len(m.order) - 1; i >= 0; i-- {
		g := m.games[m.order[i]]
		if g.Status == StatusWaiting && g.UserID1 != userID && g.UserID2 != userID {
			return g, nil
		}
	}
	return nil, ErrGameNotFound
}

// --- CreateGame ---

func TestGameService_CreateGame(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)

	g, err := svc.CreateGame(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.ID == "" {
		t.Error("expected non-empty game ID")
	}
	if g.Status != StatusWaiting {
		t.Errorf("expected status %s, got %s", StatusWaiting, g.Status)
	}
	// verify persisted
	stored, err := repo.FindByID(ctx, g.ID)
	if err != nil {
		t.Fatalf("game not found in repo: %v", err)
	}
	if stored.ID != g.ID {
		t.Error("stored game ID mismatch")
	}
}

func TestGameService_CreateGame_Errors(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	repo.createErr = errors.New("storage failure")
	svc := NewGameService(repo, slog.Default(), time.Minute)

	_, err := svc.CreateGame(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGameService_CreatePrivateGame(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)

	g, err := svc.CreatePrivateGame(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.ID == "" {
		t.Error("expected non-empty game ID")
	}
	if g.Status != StatusWaiting {
		t.Errorf("expected status %s, got %s", StatusWaiting, g.Status)
	}
	if !g.Private {
		t.Error("expected private game")
	}
}

// --- JoinGame ---

func TestGameService_JoinGame(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)
	g, _ := svc.CreateGame(ctx)

	updated, err := svc.JoinGame(ctx, g.ID, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.UserID1 != "alice" {
		t.Errorf("expected UserID1 alice, got %s", updated.UserID1)
	}
	if updated.Status != StatusWaiting {
		t.Errorf("expected status still waiting, got %s", updated.Status)
	}

	updated, err = svc.JoinGame(ctx, g.ID, "bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.UserID2 != "bob" {
		t.Errorf("expected UserID2 bob, got %s", updated.UserID2)
	}
	if updated.Status != StatusInProgress {
		t.Errorf("expected status in_progress, got %s", updated.Status)
	}
}

func TestGameService_JoinGame_Errors(t *testing.T) {
	ctx := context.Background()
	var mockRepo *mockRepo

	for _, tc := range []struct {
		name    string
		setup   func(*GameService) (gameID string, userID string)
		wantErr func() error
	}{
		{
			name: "game not found",
			setup: func(svc *GameService) (string, string) {
				return "nonexistent", "alice"
			},
			wantErr: func() error { return ErrGameNotFound },
		},
		{
			name: "game not waiting",
			setup: func(svc *GameService) (string, string) {
				g, _ := svc.CreateGame(ctx)
				svc.JoinGame(ctx, g.ID, "alice")
				svc.JoinGame(ctx, g.ID, "bob")
				return g.ID, "charlie"
			},
			wantErr: func() error { return ErrGameNotWaiting },
		},
		{
			name: "repo update fails",
			setup: func(svc *GameService) (string, string) {
				g, _ := svc.CreateGame(ctx)
				svc.JoinGame(ctx, g.ID, "alice")
				mockRepo.updateErr = errors.New("update failure")
				return g.ID, "bob"
			},
			wantErr: func() error { return mockRepo.updateErr },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo = newMockRepo()
			svc := NewGameService(mockRepo, slog.Default(), time.Minute)
			gameID, userID := tc.setup(svc)
			_, err := svc.JoinGame(ctx, gameID, userID)
			if !errors.Is(err, tc.wantErr()) {
				t.Errorf("expected %v, got %v", tc.wantErr(), err)
			}
		})
	}
}

// --- GetLatestGameForUser ---

func TestGameService_GetLatestGameForUser(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)

	// no games for user
	_, err := svc.GetLatestGameForUser(ctx, "alice")
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got %v", err)
	}

	// create some games
	g1, _ := svc.CreateGame(ctx)
	svc.JoinGame(ctx, g1.ID, "alice")
	g2, _ := svc.CreateGame(ctx)
	svc.JoinGame(ctx, g2.ID, "alice")
	svc.JoinGame(ctx, g2.ID, "bob")

	latest, err := svc.GetLatestGameForUser(ctx, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest.ID != g2.ID {
		t.Errorf("expected latest game ID %s, got %s", g2.ID, latest.ID)
	}
}

func TestGameService_GetLatestGameForUser_Errors(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	repo.getErr = errors.New("storage failure")
	svc := NewGameService(repo, slog.Default(), time.Minute)

	_, err := svc.GetLatestGameForUser(ctx, "alice")
	if !errors.Is(err, repo.getErr) {
		t.Errorf("expected %v, got %v", repo.getErr, err)
	}
}

// --- helper to create a game in progress between two users ---

func createGameInProgressHelper(svc *GameService, userID1, userID2 string) (*Game, error) {
	ctx := context.Background()
	g, err := svc.CreateGame(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := svc.JoinGame(ctx, g.ID, userID1); err != nil {
		return nil, err
	}
	if _, err := svc.JoinGame(ctx, g.ID, userID2); err != nil {
		return nil, err
	}
	return g, nil
}

// --- GetGame ---

func TestGameService_GetGame(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)
	g, _ := createGameInProgressHelper(svc, "alice", "bob")

	got, err := svc.GetGame(ctx, g.ID, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != g.ID {
		t.Error("game ID mismatch")
	}
	if got.UserID1LastSeen == nil {
		t.Error("expected UserID1LastSeen to be set after GetGame by participant")
	}
}

func TestGameService_GetGame_TouchesPresence(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)
	// Insert a fresh in-progress game with no presence recorded yet.
	g := NewGameInProgress("game-presence", "alice", "bob")
	if err := repo.Create(ctx, g); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// A non-participant read must not stamp either presence field.
	got, err := svc.GetGame(ctx, g.ID, "carol")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserID1LastSeen != nil || got.UserID2LastSeen != nil {
		t.Error("non-participant GetGame should not touch presence")
	}

	// bob reading the game stamps only bob's field.
	got, err = svc.GetGame(ctx, g.ID, "bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserID2LastSeen == nil {
		t.Error("expected UserID2LastSeen to be set after GetGame by bob")
	}
	if got.UserID1LastSeen != nil {
		t.Error("GetGame by bob should not stamp alice's presence")
	}
}

func TestGameService_GetGame_Errors(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)
	_, err := svc.GetGame(ctx, "nonexistent", "alice")
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got %v", err)
	}
}

// --- MakeMove ---

func TestGameService_MakeMove(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)
	g, _ := createGameInProgressHelper(svc, "alice", "bob")

	updated, err := svc.MakeMove(ctx, g.ID, "alice", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Board[0][0] != 1 {
		t.Error("expected mark 1 at (0,0)")
	}
	if updated.CurrentPlayerID != "bob" {
		t.Errorf("expected bob to be next, got %s", updated.CurrentPlayerID)
	}
}

func TestGameService_MakeMove_Errors(t *testing.T) {
	ctx := context.Background()
	var mockRepo *mockRepo

	for _, tc := range []struct {
		name    string
		setup   func(*GameService) (gameID string, userID string, x, y int)
		wantErr func() error
	}{
		{
			name: "game not found",
			setup: func(svc *GameService) (string, string, int, int) {
				return "nonexistent", "alice", 0, 0
			},
			wantErr: func() error { return ErrGameNotFound },
		},
		{
			name: "game move fails",
			setup: func(svc *GameService) (string, string, int, int) {
				g, _ := createGameInProgressHelper(svc, "alice", "bob")
				return g.ID, "bob", 0, 0
			},
			wantErr: func() error { return ErrNotYourTurn },
		},
		{
			name: "update fails",
			setup: func(svc *GameService) (string, string, int, int) {
				g, _ := createGameInProgressHelper(svc, "alice", "bob")
				mockRepo.updateErr = errors.New("update failure")
				return g.ID, "alice", 0, 0
			},
			wantErr: func() error { return mockRepo.updateErr },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo = newMockRepo()
			svc := NewGameService(mockRepo, slog.Default(), time.Minute)
			gameID, userID, x, y := tc.setup(svc)
			_, err := svc.MakeMove(ctx, gameID, userID, x, y)
			if !errors.Is(err, tc.wantErr()) {
				t.Errorf("expected %v, got %v", tc.wantErr(), err)
			}
		})
	}
}

// --- GiveUp ---

func TestGameService_GiveUp(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	svc := NewGameService(repo, slog.Default(), time.Minute)
	g, _ := createGameInProgressHelper(svc, "alice", "bob")

	updated, err := svc.GiveUp(ctx, g.ID, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != StatusFinished {
		t.Errorf("expected finished, got %s", updated.Status)
	}
	if updated.WinnerID != "bob" {
		t.Errorf("expected winner bob, got %s", updated.WinnerID)
	}
}

func TestGameService_GiveUp_Errors(t *testing.T) {
	ctx := context.Background()
	var mockRepo *mockRepo

	for _, tc := range []struct {
		name    string
		setup   func(*GameService) (gameID string, userID string)
		wantErr func() error
	}{
		{
			name: "game not found",
			setup: func(svc *GameService) (string, string) {
				return "nonexistent", "alice"
			},
			wantErr: func() error { return ErrGameNotFound },
		},
		{
			name: "game action fails",
			setup: func(svc *GameService) (string, string) {
				g, _ := createGameInProgressHelper(svc, "alice", "bob")
				return g.ID, "charlie"
			},
			wantErr: func() error { return ErrNotInGame },
		},
		{
			name: "update fails",
			setup: func(svc *GameService) (string, string) {
				g, _ := createGameInProgressHelper(svc, "alice", "bob")
				mockRepo.updateErr = errors.New("update failure")
				return g.ID, "alice"
			},
			wantErr: func() error { return mockRepo.updateErr },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo = newMockRepo()
			svc := NewGameService(mockRepo, slog.Default(), time.Minute)
			gameID, userID := tc.setup(svc)
			_, err := svc.GiveUp(ctx, gameID, userID)
			if !errors.Is(err, tc.wantErr()) {
				t.Errorf("expected %v, got %v", tc.wantErr(), err)
			}
		})
	}
}
