package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read(path string) Config {
	data, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return Config{}
	}

	var cfg Config

	byteStream, err := io.ReadAll(data)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return Config{}
	}
	fmt.Println(string(byteStream))

	defer data.Close()

	if err := json.Unmarshal(byteStream, &cfg); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	return cfg
}

func write(config Config) error {
	jsonData, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
		return err
	}

	path, err := GetConfigFilePath()
	if err != nil {
		log.Fatalf("Error getting file path: %v", err)
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return err
	}

	file.WriteString(string(jsonData))
	return nil
}

func GetConfigFilePath() (string, error) {
	UserHomeDir, err := os.UserHomeDir()
	path := UserHomeDir + "/" + configFileName
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "", err
	} else {
		return path, nil
	}
}

func (cfg Config) SetUser(username string) {
	cfg.Current_user_name = username
	if err := write(cfg); err != nil {
		log.Fatalf("Error changing the username: %v", err)
	}
}
