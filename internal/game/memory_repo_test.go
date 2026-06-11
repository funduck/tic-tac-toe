package game

import (
	"errors"
	"testing"
)

func TestMemoryRepo_Create(t *testing.T) {
	repo := NewMemoryRepo()
	g := NewGameInProgress("id1", "alice", "bob")
	if err := repo.Create(g); err != nil {
		t.Fatalf("Create error: %v", err)
	}
}

func TestMemoryRepo_FindByID(t *testing.T) {
	repo := NewMemoryRepo()
	g := NewGameInProgress("id1", "alice", "bob")
	if err := repo.Create(g); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	got, err := repo.FindByID("id1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != g.ID {
		t.Errorf("expected ID %s, got %s", g.ID, got.ID)
	}
}

func TestMemoryRepo_FindByID_NotFound(t *testing.T) {
	repo := NewMemoryRepo()
	_, err := repo.FindByID("nonexistent")
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got %v", err)
	}
}

func TestMemoryRepo_FindLatestForUser(t *testing.T) {
	repo := NewMemoryRepo()
	g1 := NewGameInProgress("id1", "alice", "bob")
	g2 := NewGameInProgress("id2", "charlie", "alice")
	if err := repo.Create(g1); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if err := repo.Create(g2); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	got, err := repo.FindLatestForUser("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != g2.ID {
		t.Errorf("expected latest game ID %s, got %s", g2.ID, got.ID)
	}
}

func TestMemoryRepo_FindLatestForUser_NotFound(t *testing.T) {
	repo := NewMemoryRepo()
	_, err := repo.FindLatestForUser("nonexistent")
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got %v", err)
	}
}

func TestMemoryRepo_Update(t *testing.T) {
	repo := NewMemoryRepo()
	g := NewGameInProgress("id1", "alice", "bob")
	if err := repo.Create(g); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	g.Status = StatusFinished
	if err := repo.Update(g); err != nil {
		t.Fatalf("Update error: %v", err)
	}

	got, err := repo.FindByID("id1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != StatusFinished {
		t.Errorf("expected status %s, got %s", StatusFinished, got.Status)
	}
}

func TestMemoryRepo_Update_NotFound(t *testing.T) {
	repo := NewMemoryRepo()
	g := NewGameInProgress("id1", "alice", "bob")
	err := repo.Update(g)
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got %v", err)
	}
}

func TestMemoryRepo_FindGameToJoin(t *testing.T) {
	repo := NewMemoryRepo()
	g1 := NewGame("id1")                           // waiting game
	g2 := NewGameInProgress("id2", "alice", "bob") // in_progress game
	if err := repo.Create(g1); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if err := repo.Create(g2); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	got, err := repo.FindGameToJoin("charlie")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != g1.ID {
		t.Errorf("expected to find waiting game ID %s, got %s", g1.ID, got.ID)
	}
}

func TestMemoryRepo_FindGameToJoin_NotFound(t *testing.T) {
	repo := NewMemoryRepo()
	_, err := repo.FindGameToJoin("charlie")
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got %v", err)
	}

	// also test when there are only in_progress games
	g := NewGameInProgress("id1", "alice", "bob")
	if err := repo.Create(g); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	_, err = repo.FindGameToJoin("charlie")
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got %v", err)
	}

	// and when there are only private waiting games
	g2 := NewGame("id2")
	g2.Private = true
	if err := repo.Create(g2); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	_, err = repo.FindGameToJoin("charlie")
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got %v", err)
	}
}
