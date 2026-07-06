package user

import (
	"context"
	"errors"
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

func TestMemoryUserRepo_Create_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryUserRepo()

	if err := repo.Create(ctx, &User{ID: "user3", Password: "pass"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err := repo.Create(ctx, &User{ID: "user3", Password: "other"})
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}

	found, err := repo.FindByID(ctx, "user3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Password != "pass" {
		t.Errorf("existing user was overwritten: password %q", found.Password)
	}
}
