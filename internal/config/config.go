package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gator/internal/database"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const configFileName = ".gatorconfig.json"

type Config struct { //DB connec config w JSON attachment
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name,omitempty"`
}
type State struct {
	ConfigPtr *Config
	Db        *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commandslist map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.Commandslist == nil {
		c.Commandslist = make(map[string]func(*State, Command) error)
	}
	c.Commandslist[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {

	handler, exists := c.Commandslist[cmd.Name]
	if !exists {
		return fmt.Errorf("Unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)

}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Login username required")
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
		return errors.New("Username required for registration")
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

func Read() (Config, error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		return Config{}, nil
	}
	data, err := os.ReadFile(filepath) //returns json object
	if err != nil {
		return Config{}, nil
	}
	var cfg Config //init holder container
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}
	finalString := filepath.Join(homedir, configFileName)
	return finalString, nil

}

func Write(cfg Config) error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", " ")

	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)

}
func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	return Write(*cfg)
}
