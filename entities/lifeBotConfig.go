package entities

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type LifeBotConfig struct {
	FileLocation          string
	AppToken              string
	BotToken              string
	ADPat                 string
	Channels              []Channel
	ExcludedSlackStatuses []string
	LastUpdated           time.Time
}

func (c LifeBotConfig) SaveCurrentState() {
	_, err := os.Stat(c.FileLocation)
	if err != nil {
		os.Create(c.FileLocation)
	}

	c.LastUpdated = time.Now()
	storedState, _ := json.MarshalIndent(c, "", "  ")
	_ = ioutil.WriteFile(c.FileLocation, storedState, 0644)
}

func (c LifeBotConfig) LoadState(fileLocation string) {
	file, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		os.Create(fileLocation)
	}
	config := LifeBotConfig{}

	if file != nil {
		_ = json.Unmarshal(file, &config)
	}

	c.FileLocation = fileLocation // migration step for older files - remove
	c = config
}
