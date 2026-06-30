package user

import (
	"context"
	"testing"
)

func TestMemoryUserRepo_FindByID(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryUserRepo()
	user := &User{ID: "user1", Password: "pass"}
	repo.Save(ctx, user)

	found, err := repo.FindByID(ctx, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find user, got nil")
	}
	if found.ID != "user1" {
		t.Errorf("expected ID 'user1', got '%s'", found.ID)
	}
}

func TestMemoryUserRepo_Save(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryUserRepo()
	user := &User{ID: "user2", Password: "pass"}
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByID(ctx, "user2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find user, got nil")
	}
	if found.ID != "user2" {
		t.Errorf("expected ID 'user2', got '%s'", found.ID)
	}
}
