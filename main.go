package main

import (
	"fmt"
	"gator/internal/config"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	fmt.Println("Initial config: ", cfg)

	err = cfg.SetUser("Don")
	if err != nil {
		log.Fatalf("Failed to set user: %v", err)
	}
	fmt.Println("User set successfully")

	updatedCfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read updated config: %v", err)
	}
	fmt.Println("Updated config: ", updatedCfg)
	//read the config file, set user to current user name, read config file again and print struct
}
