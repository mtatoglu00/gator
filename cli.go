package main

import (
	"context"
	"fmt"
	"log"

	"gator/internal/database"

	"github.com/google/uuid"
)

func (c *Commands) Run(s *State, cmd Command) error {
	if err := c.Cmds[cmd.Name](s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Cmds[name] = f
}

func ScrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("error trying to fetch next feed", err)
		return err
	}
	s.Db.MarkFeedFetched(context.Background(), feed.ID)

	fetchedFeed, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Println("error when fetching feed", err)
		return err
	}

	fmt.Println(fetchedFeed.Channel.Title)
	fmt.Println(fetchedFeed.Channel.Description)
	for _, v := range fetchedFeed.Channel.Item {
		params := database.CreatePostParams{
			ID:          uuid.New(),
			Title:       v.Title,
			Url:         v.Link,
			Description: v.Description,
			PublishedAt: v.PubDate,
			FeedID:      feed.ID,
		}
		err := s.Db.CreatePost(context.Background(), params)
		if err != nil {
			log.Println("error when creating post in db", err)
		}

	}

	return nil
}
