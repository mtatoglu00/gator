package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gator/internal/database"
	"gator/internal/rss"

	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("the login handler expects a single argument: The username")
	}

	usercheck, _ := Sta.Db.GetUser(context.Background(), cmd.Arguments[0])
	if usercheck.Name == "" {
		return errors.New("User does not exist!")
	}

	if err := s.Config.SetUser(cmd.Arguments[0]); err != nil {
		return err
	}
	fmt.Printf("User: %s has been set.", cmd.Arguments[0])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("the register handler expects a single argument: The username")
	}

	userCheck, err := Sta.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		log.Println("User not present. Creating user...")
	}

	if userCheck.Name != "" {
		log.Fatalln("User already present!")
		return errors.New("User already present!")
	}
	fmt.Println(userCheck.Name)

	params := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Arguments[0]}

	data, err := Sta.Db.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatalf("Error when creating a user: %v", err)
	}
	fmt.Printf("User %s created successfully!", data.Name)
	s.Config.SetUser(data.Name)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.Reset(context.Background())
	if err != nil {
		log.Fatalln("Failed to delete users")
		return err
	}
	fmt.Println("All users deleted")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		log.Fatalln("Failed to get all users")
		return err
	}
	for _, v := range users {
		if v.Name == s.Config.Current_user_name {
			fmt.Println(v.Name, "(current)")
		} else {
			fmt.Println(v.Name)
		}
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	rss.Agg()
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	user, err := s.Db.GetUser(context.Background(), s.Config.Current_user_name)
	if err != nil {
		log.Fatalf("error when fetching current user: %w", err)
		return err
	}

	if len(cmd.Arguments) < 2 {
		log.Fatal("error. Not enough argument")
		return nil
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    user.ID,
	}
	feed, err := s.Db.CreateFeed(context.Background(), params)
	if err != nil {
		log.Fatalf("error when creating feed: %w", err)
		return err
	}

	fmt.Println(feed.Name)
	fmt.Println(feed.Url)

	HandlerFollow(s, Command{Name: "follow", Arguments: []string{feed.Url}})
	if err != nil {
		log.Fatalln("error following feed", err)
		return err
	}
	return nil
}

func HandlerGetFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalln("error while fetching feeds")
		return err
	}

	for _, v := range feeds {
		user_id, err := s.Db.GetUserByID(context.Background(), v.UserID)
		if err != nil {
			log.Fatalln("error while fetching user by id")
			return err
		}
		fmt.Println(v.Name)
		fmt.Println(v.Url)
		fmt.Println(user_id.Name)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command) error {
	user, err := s.Db.GetUser(context.Background(), s.Config.Current_user_name)
	if err != nil {
		log.Fatalln("error fetching user")
		return err
	}

	feed, err := s.Db.GetFeedByURL(context.Background(), cmd.Arguments[0])
	if err != nil {
		log.Fatalln("error fetching feed by URL")
		return err
	}

	params := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		log.Fatalln("error when creating feed follow", err)
		return err
	}

	return nil
}

func HandlerFollowing(s *State, cmd Command) error {
	user, err := s.Db.GetUser(context.Background(), s.Config.Current_user_name)
	if err != nil {
		log.Fatalln("error when fetching current user", err)
		return err
	}

	followedFeeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		log.Fatalln("error when fetching users followed feeds", err)
		return err
	}

	for _, v := range followedFeeds {
		fmt.Println(v.FeedName, v.Url)
	}
	return nil
}
