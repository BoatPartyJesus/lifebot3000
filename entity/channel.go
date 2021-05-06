package entity

import (
	"errors"
	"math/rand"
	"time"
)

type Channel struct {
	ChannelName   string
	ChannelId     string
	EligibleUsers []string
	ExemptUsers   []string
	RecentUsers   []string
}

func (channel Channel) PickAWinner(meeseeks *MeeseeksSlack) (string, error) {
	var pickableNames []string

	availableUsers := meeseeks.GetAvailableUsers(channel.EligibleUsers)

	for _, user := range availableUsers {
		found := false
		for _, r := range channel.RecentUsers {
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
		return "", errors.New("no available users in channel")
	} else {
		return pickableNames[rand.Intn(len(pickableNames))], nil
	}
}

func (channel Channel) AddRecentUser(x string) []string {
	s := channel.RecentUsers
	s = append(s, x)
	if len(s) > 3 {
		s = s[1:]
	}
	return s
}

func (channel Channel) Remove(s []string, x string) []string {
	for index, i := range s {
		if i == x {
			s[index] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}

func (channel Channel) Find(s []string, x string) bool {
	for _, i := range s {
		if i == x {
			return true
		}
	}
	return false
}
