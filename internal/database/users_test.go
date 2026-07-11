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

	queries := database.New(db).WithTx(tx)
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

func TestGetUser(t *testing.T) {
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

	userId := uuid.New()

	newUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		ID:        userId,
		Name:      "peter",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	user, err := queries.GetUser(ctx, "peter")
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != userId {
		t.Fatalf("expected ID %q, got %q", userId, user.ID)
	}
	if user.Name != newUser.Name {
		t.Fatalf("expected name %q, got %q", newUser.Name, user.Name)
	}
}

func TestGetUsers(t *testing.T) {
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

	alice, err := queries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      "alice",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	bob, err := queries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      "bob",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	users, err := queries.GetUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if users[0].ID != alice.ID {
		t.Fatalf("expected first user ID %q, got %q", alice.ID, users[0].ID)
	}

	if users[1].ID != bob.ID {
		t.Fatalf("expected second user ID %q, got %q", bob.ID, users[1].ID)
	}
}

func TestDeleteAllUsers(t *testing.T) {
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

	_, err = queries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      "alice",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = queries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      "bob",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	users, err := queries.GetUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	err = queries.DeleteAllUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	users, err = queries.GetUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}
}
