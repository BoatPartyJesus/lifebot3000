package handlers

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"meeseeks/entities"
)

func ListHandler(event slackevents.EventsAPIEvent, client *slack.Client, botConfig entities.LifeBotConfig) entities.LifeBotConfig {
	fmt.Println("ListHandler")
	//do things
	// requiredList options
	// all : return all users in channel participating in the lottery
	// eligible : return all users currently pickable
	// exempt: return all users currently NOT pickable or not participating
	return botConfig
}
