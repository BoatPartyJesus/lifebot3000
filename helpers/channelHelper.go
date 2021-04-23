package channelHelper

import (
	"errors"
	"meeseeks/entities"
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

func AddRecentUser(s []string, x string) []string {
	s = append(s, x)
	if len(s) > 3 {
		s = s[1:]
	}
	return s
}

func PseudoRandomSelect(eligible []string, recent []string) (string, error) {

	var pickableNames []string

	for _, user := range eligible {
		found := false
		for _, r := range recent {
			if user == r {
				found = true
				break
			}
		}

		if !found {
			pickableNames = append(pickableNames, user)
		}
	}

	rand.Seed(time.Now().Unix())

	if len(pickableNames) == 0 {
		return "", errors.New("no eligible users in channel")
	} else {
		return pickableNames[rand.Intn(len(pickableNames))], nil
	}
}
