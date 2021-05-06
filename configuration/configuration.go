package configuration

import (
	"encoding/json"
	"io/ioutil"
	"meeseeks/entity"

	"os"
	"time"
)

func LoadConfiguration(configFile string) entity.MeeseeksConfig {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		os.Create(configFile)
	}

	config := entity.MeeseeksConfig{
		FileLocation: configFile,
	}

	if file != nil {
		_ = json.Unmarshal(file, &config)
	}

	return config
}

func SaveConfiguration(configFile string, currentConfiguration entity.MeeseeksConfig) {
	_, err := os.Stat(configFile)
	if err != nil {
		os.Create(configFile)
	}

	currentConfiguration.LastUpdated = time.Now()
	storedState, _ := json.MarshalIndent(currentConfiguration, "", "  ")
	_ = ioutil.WriteFile(configFile, storedState, 0644)
}
