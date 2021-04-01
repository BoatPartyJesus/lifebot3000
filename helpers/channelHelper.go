package channelHelper

import (
	"github.com/slack-go/slack"
	"lifebot3000/entities"
)

func RetrieveOrCreate(foundChannel slack.Channel, knownChannels []entities.Channel) []entities.Channel {

	found := false

	for _, ch := range knownChannels {
		if ch.ChannelId == foundChannel.ID {
			found = true
		}
	}

	if !found {
		knownChannels = append(knownChannels,
			entities.Channel{
				ChannelName:   foundChannel.NameNormalized,
				ChannelId:     foundChannel.ID,
				EligibleUsers: nil,
				ExemptUsers:   nil,
				RecentUsers:   nil,
			})
	}

	return knownChannels
}

func Find(s []string, x string) bool {
	for _, i := range s {
		if i == x {
			return true
		}
	}
	return false
}

func Remove(s []string, x string) []string {
	for index, i := range s {
		if i == x {
			s[index] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}
