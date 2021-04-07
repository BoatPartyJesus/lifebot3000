package channelHelper

import (
	"lifebot3000/entities"
	"math/rand"
	"time"

	"github.com/slack-go/slack"
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

func PseudoRandomSelect(eligible []string, recent []string) string {
	for i := 0; i < len(eligible); i++ {
		user := eligible[i]
		for _, r := range recent {
			if user == r {
				eligible = append(eligible[:i], eligible[i+1:]...)
				i--
				break
			}
		}
	}

	rand.Seed(time.Now().Unix())
	return eligible[rand.Intn(len(eligible))]
}
