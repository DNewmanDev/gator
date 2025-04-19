package main

import (
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

	inputCommand := os.Args

	if len(inputCommand) < 2 {
		fmt.Println("Input insufficient length, provide command")
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
