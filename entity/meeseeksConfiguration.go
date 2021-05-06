package entity

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/slack-go/slack"
)

type MeeseeksConfig struct {
	FileLocation          string
	AppToken              string
	BotToken              string
	ADPat                 string
	Channels              []Channel
	ExcludedSlackStatuses []string
	LastUpdated           time.Time
}

func (c MeeseeksConfig) SaveCurrentState() {
	_, err := os.Stat(c.FileLocation)
	if err != nil {
		os.Create(c.FileLocation)
	}

	c.LastUpdated = time.Now()
	storedState, _ := json.MarshalIndent(c, "", "  ")
	_ = ioutil.WriteFile(c.FileLocation, storedState, 0644)
}

func (c MeeseeksConfig) LoadState(fileLocation string) {
	file, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		os.Create(fileLocation)
	}
	config := MeeseeksConfig{}

	if file != nil {
		_ = json.Unmarshal(file, &config)
	}

	c.FileLocation = fileLocation // migration step for older files - remove
	c = config
}

func RetrieveOrCreateChannel(foundChannel slack.Channel, knownChannels []Channel) []Channel {

	found := false

	for _, ch := range knownChannels {
		if ch.ChannelId == foundChannel.ID {
			found = true
		}
	}

	if !found {
		knownChannels = append(knownChannels,
			Channel{
				ChannelName:   foundChannel.NameNormalized,
				ChannelId:     foundChannel.ID,
				EligibleUsers: nil,
				ExemptUsers:   nil,
				RecentUsers:   nil,
			})
	}

	return knownChannels
}
