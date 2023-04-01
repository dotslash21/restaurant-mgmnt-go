package configs

import (
	"encoding/json"
	"os"
)

func LoadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}

	defer file.Close()

	config := make(map[string]string)
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return err
	}

	for key, value := range config {
		os.Setenv(key, value)
	}

	return nil
}
