package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)


type Config struct {
	BaseUrl string `json:"base_url"`
	CSRFToken string `json:"csrf_token"`
	Cookie string `json:"cookie"`
}

func Load() (*Config, error){
	configContent, err := os.ReadFile("./config/config.json")

	if err != nil {
		fmt.Println("Error reading config file", err)
		return nil, err
	}

	decoder := json.NewDecoder( bytes.NewReader(configContent))
	decoder.DisallowUnknownFields()

	var cfg Config

	if err := decoder.Decode(&cfg); err != nil {
		fmt.Println("Error decoding config", err)
		return nil, err
	}

	return &cfg,nil
}
