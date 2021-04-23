package channelHelper

import (
	"github.com/slack-go/slack"
	"log"
	"meeseeks/entities"
	"os"
)

type MeeseeksSlack struct {
	Slack *slack.Client
	BotConfig *entities.LifeBotConfig
}

func (meeseeks MeeseeksSlack) New(botConfig *entities.LifeBotConfig) {
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
		userProfile, _ := meeseeks.Slack.GetUserProfile(user, true)
		if !Find(meeseeks.BotConfig.ExcludedSlackStatuses, userProfile.StatusText) {
			availableUsers = append(availableUsers, user)
		}
	}

	return availableUsers
}