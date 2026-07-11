package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rvanhoepen/gator/internal/database"
)

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		tx.Rollback()
	})

	queries := database.New(tx)
	ctx := context.Background()

	user, err := queries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      "alice",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	if user.Name != "alice" {
		t.Fatalf("expected 'alice', got '%s'", user.Name)
	}
}
