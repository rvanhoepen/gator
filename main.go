package main

import (
	"fmt"
	"io"
	"os"

	"github.com/rvanhoepen/gator/internal/config"
)

type state struct {
	input  io.Reader
	output io.Writer
	cfg    *config.Config
}

func newState(input io.Reader, output io.Writer, cfg *config.Config) *state {
	state := &state{
		input:  input,
		output: output,
		cfg:    cfg,
	}

	return state
}

func (s *state) run(c commands) {
	input := os.Args
	if len(input) < 2 {
		fmt.Fprintln(os.Stderr, "no command provided")
		os.Exit(1)
	}
	command := command{
		name: input[1],
		args: input[2:],
	}
	if err := c.run(s, command); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	state := newState(os.Stdin, os.Stdout, &cfg)
	commands := newCommands()

	state.run(commands)
}
