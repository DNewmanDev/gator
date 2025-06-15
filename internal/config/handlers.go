package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/database"
	"log"
	"os"
	"strings"
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
	//neds protection logic against incorrect user input
	time_between_reqs := cmd.Args[0]
	parsedDuration, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return err
	}
	fmt.Println("Collecting feeds every \n", parsedDuration)

	ticker := time.NewTicker(parsedDuration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

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

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("URL required")
	}
	targeturl := cmd.Args[0]

	params := database.DeleteFeedFollowParams{UserID: user.ID, Url: targeturl}
	err := s.Db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to grab delete for , error: %v", err)
	}
	fmt.Printf("Feed with URL %s unfollowed", targeturl)
	return nil
}

func scrapeFeeds(s *State) {

	//get next feed to fetch
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Failed to grab next feed, err: %v", err)
		return
	}
	//mark it as fetched
	log.Println("Fetching feed \n", nextFeed)
	scrapeFeed(s.Db, nextFeed)

}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("failed to mark next feed, err: %v", err)
		return
	}
	returnedFeed, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("failed to get feed, err: %v", err)
		return
	}
	for _, object := range returnedFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, object.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     object.Title,
			Description: sql.NullString{
				String: object.Description,
				Valid:  true,
			},
			Url:         object.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(returnedFeed.Channel.Item))
}
