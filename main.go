package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	_ "github.com/lib/pq"
	"github.com/rvanhoepen/gator/internal/config"
	"github.com/rvanhoepen/gator/internal/database"
)

type state struct {
	input  io.Reader
	output io.Writer
	cfg    *config.Config
	db     *database.Queries
}

func newState(input io.Reader, output io.Writer, cfg *config.Config, db *database.Queries) *state {
	state := &state{
		input:  input,
		output: output,
		cfg:    cfg,
		db:     db,
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
		fmt.Fprintln(os.Stderr, "config file is missing")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not open database")
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "could not close database:", err)
		}
	}()

	if err := db.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, "could not connect to database")
		os.Exit(1)
	}

	dbQueries := database.New(db)

	state := newState(os.Stdin, os.Stdout, &cfg, dbQueries)
	commands := newCommands()

	state.run(commands)
}
