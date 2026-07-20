package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rvanhoepen/gator/internal/database"
	"github.com/rvanhoepen/gator/internal/feed"
)

func scrapeFeeds(s *state) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(s.output, "[INFO] Fetching from: %s...\n", nextFeed.Url)

	feed, err := feed.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		fmt.Fprintf(s.output, "[ERROR] Failed to fetch: %v\n", err)
		return err
	}

	postsCount := 0

	for _, item := range feed.Channel.Item {
		title := strings.TrimSpace(item.Title)
		url := strings.TrimSpace(item.Link)
		descriptionText := strings.TrimSpace(item.Description)

		description := sql.NullString{
			String: descriptionText,
			Valid:  descriptionText != "",
		}

		publishedAt := sql.NullTime{}
		if item.PubDate != "" {
			t, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err == nil {
				publishedAt = sql.NullTime{
					Time:  t,
					Valid: true,
				}
			}
		}

		_, err := s.db.CreatePost(
			ctx,
			database.CreatePostParams{
				ID:          uuid.New(),
				Title:       title,
				Url:         url,
				Description: description,
				PublishedAt: publishedAt,
				FeedID:      nextFeed.ID,
			},
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue // duplicate URL with ON CONFLICT DO NOTHING
			}
			return err
		}
		postsCount++
	}

	fmt.Fprintf(s.output, "[SUCCESS] Inserted %d posts.\n", postsCount)

	return nil
}
