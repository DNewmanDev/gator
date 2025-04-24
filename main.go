package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	state := &config.State{ConfigPtr: &cfg}
	state.ConfigPtr.DBURL = "postgres://postgres:postgres@localhost:5432/gator"
	db, err := sql.Open("postgres", state.ConfigPtr.DBURL)
	if err != nil {
		log.Fatalf("Failed to open postgres connection, err: %v\n", err)
	}
	dbQueries := database.New(db)
	state.Db = dbQueries
	commandsList := &config.Commands{}
	commandsList.Register("login", config.HandlerLogin)
	commandsList.Register("register", config.HandlerRegister)
	commandsList.Register("reset", config.HandlerReset)
	commandsList.Register("users", config.HandlerList)
	commandsList.Register("agg", config.Agg)
	commandsList.Register("addfeed", middlewareLoggedIn(config.AddFeed))
	commandsList.Register("feeds", config.HandlerFeedsDisplay)
	commandsList.Register("follow", middlewareLoggedIn(config.HandlerFollow))
	commandsList.Register("following", middlewareLoggedIn(config.HandlerFollowing))

	inputCommand := os.Args

	if len(inputCommand) < 2 {
		fmt.Println("input insufficient length, provide command")
		os.Exit(1)
	}
	cmdName := inputCommand[1]
	cmdArgs := inputCommand[2:]
	cmd := config.Command{
		Name: cmdName,
		Args: cmdArgs,
	}
	if err := commandsList.Run(state, cmd); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if err := config.Write(*state.ConfigPtr); err != nil {
		fmt.Printf("Error writing config %v\n", err)
	}
}

func middlewareLoggedIn(handler func(s *config.State, cmd config.Command, user database.User) error) func(*config.State, config.Command) error {
	return func(s *config.State, cmd config.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.ConfigPtr.CurrentUserName)
		if err != nil {
			log.Fatal(err)
		}

		return handler(s, cmd, user)
	}
}
