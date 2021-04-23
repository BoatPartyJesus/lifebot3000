package handlers

import (
	"fmt"
	"meeseeks/entities"
	channelHelper "meeseeks/helpers"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
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
				requiredChannel.RecentUsers = channelHelper.Remove(requiredChannel.RecentUsers, event.User)
				requiredChannel.ExemptUsers = channelHelper.Remove(requiredChannel.ExemptUsers, event.User)
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
