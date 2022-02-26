package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
)

type Config struct {
	APIPort        string `json:"apiPort"`
	AddCORSHeaders bool   `json:"addCORSHeaders"`
}

func LoadConfig(configFiles fs.FS) (Config, error) {
	var config Config
	environment := getEnv("ENVIRONMENT", "development")

	file, err := configFiles.Open("config/" + environment + ".json")
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file for environment %s; %w", environment, err)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file for environment %s; %w", environment, err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config file for environment %s; %w", environment, err)
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
