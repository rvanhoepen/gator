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
	"github.com/rvanhoepen/gator/internal/feed"
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
			"login":     handlerLogin,
			"register":  handlerRegister,
			"reset":     handlerReset,
			"users":     handleUsers,
			"agg":       handleAgg,
			"addfeed":   middlewareLoggedIn(handleAddFeed),
			"feeds":     handleListFeeds,
			"follow":    middlewareLoggedIn(handleFollow),
			"following": middlewareLoggedIn(handleListFollowing),
			"unfollow":  middlewareLoggedIn(handleUnfollow),
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

func handleAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("agg expects 1 argument; the time between requests")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Fprintf(s.output, "Collecting feeds every %s\n", cmd.args[0])

	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Fprintln(s.output, err)
		}
	}
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("addfeed expects two arguments: name and url")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := feed.FetchFeed(ctx, url)
	if err != nil {
		return err
	}

	newFeed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Constraint == "feeds_url_key" {
			return fmt.Errorf("feed with url: %q already exists", url)
		}
		return err
	}

	fmt.Fprintf(s.output, "Feed created successfully:\n")
	fmt.Fprintf(s.output, "Name: %s\n", newFeed.Name)
	fmt.Fprintf(s.output, "URL: %s\n", newFeed.Url)

	return nil
}

func handleListFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("feeds expects no arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Fprintln(s.output, "no feeds registered yet")
		return nil
	}

	fmt.Fprint(s.output, "Listing feeds:\n")
	fmt.Fprint(s.output, "-----------------------------------\n")

	for i, feed := range feeds {
		fmt.Fprintf(s.output, "#%d: %s\n", i+1, feed.Name)
		fmt.Fprintf(s.output, "    %s\n", feed.Url)
		fmt.Fprintf(s.output, "    Created by: %s\n", feed.CreatedBy)
		fmt.Fprint(s.output, "-----------------------------------\n")
	}

	return nil
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow expects one argument: url")
	}

	url := cmd.args[0]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}

	follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Constraint == "feed_follows_user_id_feed_id_key" {
			return fmt.Errorf("you are already following this feed")
		}
		return err
	}

	fmt.Fprintf(s.output, "Successfully followed: %s\n", follow.FeedName)

	return nil
}

func handleListFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following expects no arguments")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	if len(follows) == 0 {
		fmt.Fprintln(s.output, "no feeds followed yet")
		return nil
	}

	fmt.Fprint(s.output, "Listing follows:\n")
	fmt.Fprint(s.output, "-----------------------------------\n")

	for i, follow := range follows {
		fmt.Fprintf(s.output, "#%d: %s\n", i+1, follow.FeedName)
		fmt.Fprintf(s.output, "    %s\n", follow.FeedUrl)
		fmt.Fprint(s.output, "-----------------------------------\n")
	}

	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("unfollow expects one argument: the url of the feed")
	}

	url := cmd.args[0]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(s.output, "Unfollowed %s\n", feed.Url)
	return nil
}
