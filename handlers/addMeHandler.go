package handlers

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func AddMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client) {
	fmt.Println("AddMeHandler")
	event := ev.InnerEvent.Data.(*slackevents.AppMentionEvent)

	_, _, err := client.PostMessage(event.Channel, slack.MsgOptionText("No.", false))
	if err != nil {
		fmt.Println(":(")
	}
}
