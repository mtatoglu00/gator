package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gator/internal/database"

	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("the login handler expects a single argument: The username")
	}

	usercheck, _ := Sta.Db.GetUser(context.Background(), cmd.Arguments[0])
	if usercheck == "" {
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

	if userCheck != "" {
		log.Fatalln("User already present!")
		return errors.New("User already present!")
	}
	fmt.Println(userCheck)

	params := database.CreateUserParams{ID: uuid.NullUUID{UUID: uuid.New()}, CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Arguments[0]}

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
