package configuration

import (
	"awesomeProject/entities"
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadConfiguration(configFile string) entities.LifeBotConfig {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		os.Create(configFile)
	}
	config := entities.LifeBotConfig{}

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

	storedState, _ := json.MarshalIndent(currentConfiguration, "", "  ")
	_ = ioutil.WriteFile(configFile, storedState, 0644)
}
