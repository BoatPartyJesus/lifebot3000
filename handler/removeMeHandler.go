package handler

import (
	"fmt"
	"meeseeks/entity"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func RemoveMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig) entity.MeeseeksConfig {
	fmt.Println("RemoveMeHandler")
	event := ev.InnerEvent.Data.(*slackevents.AppMentionEvent)

	for index, ch := range botConfig.Channels {
		if ch.ChannelId == event.Channel {
			requiredChannel := &botConfig.Channels[index]

			var message string

			if ch.Find(requiredChannel.EligibleUsers, event.User) {
				message = fmt.Sprintf("OK, I've removed <@%s> from %s", event.User, requiredChannel.ChannelName)
				requiredChannel.EligibleUsers = ch.Remove(requiredChannel.EligibleUsers, event.User)
				requiredChannel.RecentUsers = ch.Remove(requiredChannel.RecentUsers, event.User)
				requiredChannel.ExemptUsers = ch.Remove(requiredChannel.ExemptUsers, event.User)
			} else {
				message = fmt.Sprintf("You weren't in %s...", requiredChannel.ChannelName)
			}

			_, _, err := client.PostMessage(event.Channel, slack.MsgOptionText(message, false))
			if err != nil {
				fmt.Println(":(")
			}
		}
	}

	botConfig.SaveCurrentState()

	return botConfig
}
