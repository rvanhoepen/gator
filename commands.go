package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rvanhoepen/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registry map[string]func(*state, command) error
}

func newCommands() commands {
	return commands{
		registry: map[string]func(*state, command) error{
			"login":    handlerLogin,
			"register": handlerRegister,
			"reset":    handlerReset,
			"users":    handleUsers,
		},
	}
}

func (c *commands) run(s *state, cmd command) error {
	callback, ok := c.registry[cmd.name]
	if !ok {
		return fmt.Errorf("command not found")
	}
	return callback(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login expects the username as the only argument")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user %q does not exist", username)
		}
		return err
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Fprintf(s.output, "User `%s` has been logged in.\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("register expects the username as the only argument")
	}

	username := cmd.args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Constraint == "users_name_key" {
			return fmt.Errorf("user %q already exists", username)
		}
		return err
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Fprintf(s.output, "User `%s` created successfully.\n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset expects no arguments")
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(""); err != nil {
		return err
	}

	fmt.Fprintf(s.output, "Users cleared successfully.\n")

	return nil
}

func handleUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("users expects no arguments")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Fprintln(s.output, "no users registered yet")
		return nil
	}

	fmt.Fprintln(s.output, "available users:")
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Fprintf(s.output, "* %s (current)\n", user.Name)
		} else {
			fmt.Fprintf(s.output, "* %s\n", user.Name)
		}
	}

	return nil
}
