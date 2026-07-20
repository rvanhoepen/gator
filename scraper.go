package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rvanhoepen/gator/internal/feed"
)

func scrapeFeeds(s *state) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	feed, err := feed.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	fmt.Fprintln(s.output, "New posts imported:")

	for i, item := range feed.Channel.Item {
		fmt.Fprintf(s.output, "%d) %s\n", i+1, item.Title)
	}

	return nil
}
