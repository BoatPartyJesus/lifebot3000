package entity

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
)

type MeeseeksSlack struct {
	Slack     *slack.Client
	BotConfig *MeeseeksConfig
}

func (meeseeks *MeeseeksSlack) New(botConfig *MeeseeksConfig) {
	meeseeks.BotConfig = botConfig
	meeseeks.Slack = slack.New(
		botConfig.BotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(botConfig.AppToken))
}

func (meeseeks *MeeseeksSlack) GetAvailableUsers(users []string) []string {
	var availableUsers []string

	for _, user := range users {
		getUserParams := slack.GetUserProfileParameters{
			UserID:        user,
			IncludeLabels: false,
		}
		userProfile, err := meeseeks.Slack.GetUserProfile(&getUserParams)

		if err != nil {
			fmt.Println(err)
		}

		if !arraySeek(meeseeks.BotConfig.ExcludedSlackStatuses, userProfile.StatusText) {
			availableUsers = append(availableUsers, user)
		}
	}

	return availableUsers
}

func arraySeek(s []string, x string) bool {
	for _, i := range s {
		if i == x {
			return true
		}
	}
	return false
}
