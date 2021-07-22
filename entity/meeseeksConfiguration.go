package entity

import (
	"encoding/json"
	"meeseeks/util"
	"time"

	"github.com/slack-go/slack"
)

type IMeeseeksConfiguration interface {
	SaveCurrentState()
	LoadConfigurationStateFromFile(fileLocation string) MeeseeksConfig
	RetrieveOrCreateChannel(foundChannel slack.Channel, knownChannels []Channel) []Channel
}

type MeeseeksConfig struct {
	FileLocation          string
	AppToken              string
	BotToken              string
	ADPat                 string
	Channels              []Channel
	ExcludedSlackStatuses []string
	LastUpdated           time.Time
}

func (c MeeseeksConfig) SaveCurrentState(fu util.FileUtility) {
	if !fu.DoesFileExist(c.FileLocation) {
		fu.CreateFile(c.FileLocation)
	}

	c.LastUpdated = time.Now()
	storedState, _ := json.MarshalIndent(c, "", "  ")

	fu.WriteFile(c.FileLocation, storedState)
}

func LoadConfigurationStateFromFile(fu util.FileUtility, fileLocation string) MeeseeksConfig {
	if !fu.DoesFileExist(fileLocation) {
		fu.CreateFile(fileLocation)
	}

	file, _ := fu.ReadFile(fileLocation)

	config := MeeseeksConfig{}

	if file != nil {
		_ = json.Unmarshal(file, &config)
	}

	config.FileLocation = fileLocation // migration step for older files - remove
	return config
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
				RecentUsers:   nil,
			})
	}

	return knownChannels
}
