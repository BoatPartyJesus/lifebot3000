package entity

import (
	"errors"
	"math/rand"
	"time"
)

type IChannel interface {
	PickAWinner(meeseeks IMeeseeksSlack) (string, error)
	IsEligibleUser(x string) bool
	AddEligibleUser(x string) []string
	RemoveEligibleUser(x string) []string
	AddRecentUser(x string) []string
}

type Channel struct {
	ChannelName   string
	ChannelId     string
	EligibleUsers []string
	RecentUsers   []string
}

func (channel Channel) PickAWinner(meeseeks IMeeseeksSlack) (string, error) {
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

func (channel Channel) IsEligibleUser(x string) bool {
	s := channel.EligibleUsers
	return find(s, x)
}

func (channel Channel) AddEligibleUser(x string) []string {
	s := channel.EligibleUsers
	if !find(s, x) {
		s = append(s, x)
	}
	return s
}

func (channel Channel) RemoveEligibleUser(x string) ([]string, []string) {
	e := channel.EligibleUsers
	r := channel.RecentUsers

	if find(e, x) {
		e = remove(e, x)
		r = remove(r, x)
	}

	return e, r
}

func (channel Channel) AddRecentUser(x string) []string {
	s := channel.RecentUsers
	if !find(s, x) {
		s = append(s, x)
	}
	if len(s) > 3 {
		s = s[1:]
	}
	return s
}

func remove(s []string, x string) []string {
	for index, i := range s {
		if i == x {
			s[index] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}

func find(s []string, x string) bool {
	for _, i := range s {
		if i == x {
			return true
		}
	}
	return false
}
