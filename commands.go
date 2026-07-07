package main

import (
	"fmt"
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
			"login": handlerLogin,
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
	if len(cmd.args) == 0 {
		return fmt.Errorf("login expects the username as the only argument")
	}

	username := cmd.args[0]
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Fprintf(s.output, "User `%s` has been logged in.\n", username)

	return nil
}
