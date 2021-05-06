package handler

import (
	"fmt"
	"meeseeks/entity"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func ChannelAddMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig) entity.MeeseeksConfig {
	fmt.Println("AddMeHandler")
	event := ev.InnerEvent.Data.(*slackevents.AppMentionEvent)

	for index, ch := range botConfig.Channels {
		if ch.ChannelId == event.Channel {
			requiredChannel := &botConfig.Channels[index]

			var message string

			if ch.Find(requiredChannel.EligibleUsers, event.User) {
				message = fmt.Sprintf("You are already in %s", requiredChannel.ChannelName)
			} else {
				requiredChannel.EligibleUsers = append(requiredChannel.EligibleUsers, event.User)
				message = fmt.Sprintf("OK, I'll add <@%s> to %s", event.User, requiredChannel.ChannelName)
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

func MessageAddMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig) entity.MeeseeksConfig {
	fmt.Println("AddMeHandler")
	event := ev.InnerEvent.Data.(*slackevents.MessageEvent)

	_, _, _ = client.PostMessage(event.Channel, slack.MsgOptionText("Err... Ok. You win. Congratulations, I guess... :tada:", false))

	return botConfig
}
