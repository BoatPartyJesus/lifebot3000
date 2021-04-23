package configuration

import (
	"encoding/json"
	"io/ioutil"
	"meeseeks/entities"
	"os"
	"time"
)

func LoadConfiguration(configFile string) entities.LifeBotConfig {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		os.Create(configFile)
	}

	config := entities.LifeBotConfig{
		FileLocation: configFile,
	}

	if file != nil {
		_ = json.Unmarshal(file, &config)
	}

	return config
}

func SaveConfiguration(configFile string, currentConfiguration entities.LifeBotConfig) {
	_, err := os.Stat(configFile)
	if err != nil {
		os.Create(configFile)
	}

	currentConfiguration.LastUpdated = time.Now()
	storedState, _ := json.MarshalIndent(currentConfiguration, "", "  ")
	_ = ioutil.WriteFile(configFile, storedState, 0644)
}
