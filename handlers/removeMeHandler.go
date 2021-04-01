package handlers

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"lifebot3000/entities"
	channelHelper "lifebot3000/helpers"
)

func RemoveMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entities.LifeBotConfig) entities.LifeBotConfig {
	fmt.Println("RemoveMeHandler")
	event := ev.InnerEvent.Data.(*slackevents.AppMentionEvent)

	for index, ch := range botConfig.Channels {
		if ch.ChannelId == event.Channel {
			requiredChannel := &botConfig.Channels[index]

			var message string

			if channelHelper.Find(requiredChannel.EligibleUsers, event.User) {
				message = fmt.Sprintf("OK, I've removed <@%s> from %s", event.User, requiredChannel.ChannelName)
				requiredChannel.EligibleUsers = channelHelper.Remove(requiredChannel.EligibleUsers, event.User)
			} else {
				message = fmt.Sprintf("You weren't in %s...", requiredChannel.ChannelName)
			}

			_, _, err := client.PostMessage(event.Channel, slack.MsgOptionText(message, false))
			if err != nil {
				fmt.Println(":(")
			}
		}
	}
	return botConfig
}
