package handler

import (
	"fmt"
	"meeseeks/entity"
	"meeseeks/util"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type IAddMeHandler interface {
	AddMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig
}

// TODO: finish mocking implementation
// type IMessageClient interface{
// 	PostMessage() (string, string, error)
// }

func AddMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig {
	event := ev.InnerEvent.Data.(*slackevents.AppMentionEvent)

	user := namedUser
	if user == "" {
		user = event.User
	}

	for index, ch := range botConfig.Channels {
		if ch.ChannelId == event.Channel {
			requiredChannel := &botConfig.Channels[index]

			var message string

			if ch.IsEligibleUser(user) {
				if namedUser == "" {
					message = fmt.Sprintf("You are already in %s", requiredChannel.ChannelName)
				} else {
					message = fmt.Sprintf("<@%s> is already in %s", user, requiredChannel.ChannelName)
				}
			} else {
				requiredChannel.EligibleUsers = ch.AddEligibleUser(user)
				message = fmt.Sprintf("OK, I'll add <@%s> to %s", user, requiredChannel.ChannelName)
			}

			_, _, err := client.PostMessage(event.Channel, slack.MsgOptionText(message, false))
			if err != nil {
				fmt.Println(":(")
			}
		}
	}

	botConfig.SaveCurrentState(util.OsFileUtility{})

	return botConfig
}

// func InteractionAddMeHandler(ev slack.InteractionCallback, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig {
// 	user := namedUser
// 	if user == "" {
// 		user = ev.User.ID
// 	}

// 	for index, ch := range botConfig.Channels {
// 		if ch.ChannelId == ev.Channel.ID {
// 			requiredChannel := &botConfig.Channels[index]

// 			var message string

// 			if ch.IsEligibleUser(user) {
// 				if namedUser == "" {
// 					message = fmt.Sprintf("You are already in %s", requiredChannel.ChannelName)
// 				} else {
// 					message = fmt.Sprintf("<@%s> is already in %s", user, requiredChannel.ChannelName)
// 				}
// 			} else {
// 				requiredChannel.EligibleUsers = ch.AddEligibleUser(user)
// 				message = fmt.Sprintf("OK, I'll add <@%s> to %s", user, requiredChannel.ChannelName)
// 			}

// 			_, _, err := client.PostMessage(ev.Channel.ID, slack.MsgOptionText(message, false))
// 			if err != nil {
// 				fmt.Println(":(")
// 			}
// 		}
// 	}

// 	botConfig.SaveCurrentState()

// 	return botConfig
// }
