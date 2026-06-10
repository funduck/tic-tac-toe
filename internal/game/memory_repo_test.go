package game

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestMemoryRepo_CreateConcurrently(t *testing.T) {
	nWorkers := 10

	repo := NewMemoryRepo()

	wg := sync.WaitGroup{}
	wg.Add(nWorkers)

	errors := make(chan error, nWorkers)

	for i := range nWorkers {
		go func(i int) {
			defer wg.Done()
			g := NewGameInProgress(
				fmt.Sprintf("game%d", i),
				fmt.Sprintf("user%d", i),
				fmt.Sprintf("user%d", i+1),
			)
			if err := repo.Create(g); err != nil {
				errors <- fmt.Errorf("Create error for game%d: %v", i, err)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}

	if len(repo.gamesByID) != nWorkers {
		t.Errorf("expected %d games, got %d", nWorkers, len(repo.gamesByID))
	}
}

func TestMemoryRepo_FindByIDConcurrently(t *testing.T) {
	repo := NewMemoryRepo()
	g := NewGameInProgress("id1", "alice", "bob")
	if err := repo.Create(g); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	nWorkers := 5
	wg := sync.WaitGroup{}
	wg.Add(nWorkers)

	errors := make(chan error, nWorkers)

	for range nWorkers {
		go func() {
			defer wg.Done()

			got, err := repo.FindByID("id1")
			if err != nil {
				errors <- fmt.Errorf("FindByID error: %v", err)
				return
			}
			if got.ID != g.ID {
				errors <- fmt.Errorf("expected ID %s, got %s", g.ID, got.ID)
			}

		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
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
