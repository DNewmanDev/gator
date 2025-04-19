package config

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/database"
	"os"
	"time"

	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("login username required")
	}

	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("User with name %s does not exist\n", cmd.Args[0])

		os.Exit(1)
	}

	s.ConfigPtr.CurrentUserName = cmd.Args[0]
	fmt.Printf("Current user set to %s\n", s.ConfigPtr.CurrentUserName)

	return nil
}

func HandlerRegister(s *State, cmd Command) error {

	if len(cmd.Args) == 0 {
		return errors.New("username required for registration")
	}

	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		fmt.Printf("User with name '%s' already exists\n", cmd.Args[0])
		os.Exit(1)
	}
	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	newUser, err := s.Db.CreateUser(context.Background(), user)
	if err != nil {
		println("User creation failed, user already exists")
		os.Exit(1)
	}
	s.ConfigPtr.CurrentUserName = newUser.Name
	fmt.Printf("User creation successful, Name: %s ID: %v Time: %v\n", user.Name, user.ID, user.CreatedAt)
	return nil
}

func HandlerReset(s *State, cmd Command) error {

	s.Db.ResetTable(context.Background())

	return nil
}
func HandlerList(s *State, cmd Command) error {

	users, err := s.Db.ListUsers(context.Background())
	if err != nil {
		fmt.Println("Failed to find users for listing")
		os.Exit(1)
	}
	currentUser := s.ConfigPtr.CurrentUserName
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func Agg(s *State, cmd Command) error {
	feed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("THIS IS THE FEEEEED RESPONSE FROM AGG! :d %v", feed)
	fmt.Println(feed.Channel.Description)
	fmt.Println("Channel title: \n", feed.Channel.Title)
	fmt.Println("Channel description: \n", feed.Channel.Description)
	fmt.Println("Channel link: \n", feed.Channel.Link)

	fmt.Println("Channel contains the following items: ")
	for _, item := range feed.Channel.Item {
		fmt.Println("--Item title: ", item.Title)
		fmt.Println("--Item description: ", item.Description)
		fmt.Println("--Item link: ", item.Link)
		fmt.Println("--Item published: ", item.PubDate)
	}
	return nil
}
