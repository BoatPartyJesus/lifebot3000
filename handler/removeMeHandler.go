package handler

import (
	"fmt"
	"meeseeks/entity"
	"meeseeks/util"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type IRemoveMeHandler interface {
	RemoveMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig
}

func RemoveMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig {
	event := ev.InnerEvent.Data.(*slackevents.AppMentionEvent)

	user := namedUser
	if user == "" {
		user = event.User
	}

	// for index, ch := range botConfig.Channels {
	// 	if ch.ChannelId == event.Channel {
	// 		requiredChannel := &botConfig.Channels[index]

	// 		var message string

	// 		if ch.Find(requiredChannel.EligibleUsers, user) {
	// 			message = fmt.Sprintf("OK, I've removed <@%s> from %s", user, requiredChannel.ChannelName)
	// 			requiredChannel.EligibleUsers = ch.Remove(requiredChannel.EligibleUsers, user)
	// 			requiredChannel.RecentUsers = ch.Remove(requiredChannel.RecentUsers, user)
	// 		} else {
	// 			if namedUser == "" {
	// 				message = fmt.Sprintf("You weren't in %s...", requiredChannel.ChannelName)
	// 			} else {
	// 				message = fmt.Sprintf("<@%s> wasn't in %s...", user, requiredChannel.ChannelName)
	// 			}
	// 		}

	// 		_, _, err := client.PostMessage(event.Channel, slack.MsgOptionText(message, false))
	// 		if err != nil {
	// 			fmt.Println(":(")
	// 		}
	// 	}
	// }

	for index, ch := range botConfig.Channels {
		if ch.ChannelId == event.Channel {
			requiredChannel := &botConfig.Channels[index]

			var message string

			if ch.IsEligibleUser(user) {
				message = fmt.Sprintf("OK, I've removed <@%s> from %s", user, requiredChannel.ChannelName)
				requiredChannel.EligibleUsers, requiredChannel.RecentUsers = ch.RemoveEligibleUser(user)
			} else {
				if namedUser == "" {
					message = fmt.Sprintf("You weren't in %s...", requiredChannel.ChannelName)
				} else {
					message = fmt.Sprintf("<@%s> wasn't in %s...", user, requiredChannel.ChannelName)
				}
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
