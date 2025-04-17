package main

import (
	"fmt"
	"gator/internal/config"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	state := &config.State{ConfigPtr: &cfg}
	commandsList := &config.Commands{}
	commandsList.Register("login", config.HandlerLogin)

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
