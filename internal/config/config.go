package config

import (
	"encoding/json"
	"fmt"
	"gator/internal/database"
	"os"
	"path/filepath"
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
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)

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
