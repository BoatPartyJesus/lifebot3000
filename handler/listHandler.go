package handler

import (
	"fmt"
	"meeseeks/entity"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type IListHandler interface {
	ListHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig
}

func ListHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig {
	event := ev.InnerEvent.Data.(*slackevents.AppMentionEvent)

	for index, ch := range botConfig.Channels {
		if ch.ChannelId == event.Channel {
			requiredChannel := &botConfig.Channels[index]

			message := "All users selectable for random scrum master: \n"

			for _, user := range requiredChannel.EligibleUsers {
				message += fmt.Sprintf("<@%s>\n", user)
			}

			message += "\nUsers recently selected for random scrum master: \n"
			for _, user := range requiredChannel.RecentUsers {
				message += fmt.Sprintf("<@%s>\n", user)
			}

			_, _, err := client.PostMessage(event.Channel, slack.MsgOptionText(message, false))
			if err != nil {
				fmt.Println(":(")
			}
		}
	}
	//do things
	// requiredList options
	// all : return all users in channel participating in the lottery
	// eligible : return all users currently pickable
	// exempt: return all users currently NOT pickable or not participating
	return botConfig
}
