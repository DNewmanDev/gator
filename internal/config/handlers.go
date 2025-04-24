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

func AddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("requires feed name and URL")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{Name: name, Url: url, UserID: uuid.NullUUID{UUID: user.ID, Valid: true}})

	if err != nil {
		return fmt.Errorf("failed to create feed: %v", err)
	}
	fmt.Printf("ID=%s, Name=%s, URL=%s\n", feed.ID, feed.Name, feed.Url)

	followParams := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	feedFollows, err := s.Db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return fmt.Errorf("failed to create feed follow record: %v", err)
	}
	if len(feedFollows) > 0 {
		feedFollow := feedFollows[0]
		fmt.Printf("User %s is now following feed %s\n", feedFollow.UserName, feedFollow.FeedName)
	}

	return nil
}

func HandlerFeedsDisplay(s *State, cmd Command) error {
	feeds, err := s.Db.ListFeedsWithUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Feed Name: %v\n", feed.FeedName)
		fmt.Printf("Feed URL: %v\n", feed.FeedsUrl)
		fmt.Printf("Feed Adder: %v\n", feed.UserName)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("URL required")
	}
	//need a query to grab feed by URL
	url := cmd.Args[0]

	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to retrieve feed, error: %v", err)
	}
	params := database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID}

	feedFollows, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create feed follows entry, error: %v", err)
	}
	if len(feedFollows) > 0 {
		feedFollow := feedFollows[0]
		fmt.Printf("User %s is now following feed %s\n", feedFollow.UserName, feedFollow.FeedName)
	}
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to grab feedfollows for currentf user, error: %v", err)
	}
	if len(feedFollows) == 0 {
		fmt.Println("Follow command success, no feeds followed")
		return nil
	}
	fmt.Println("Followed feeds: ")
	for i, ff := range feedFollows {
		fmt.Printf("%d. %s\n", i+1, ff.FeedName)
	}
	return nil
}
