package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	LogFile    string
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return nil, err
	}
	config.ServerPort = serverPort

	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		return nil, fmt.Errorf("LOG_FILE environment variable not set")
	}
	config.LogFile = logFile

	return config, nil
}
