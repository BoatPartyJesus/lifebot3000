package handler

import (
	"meeseeks/entity"
	"meeseeks/util"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type IResetRecentHandler interface {
	ResetMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig
}

func ResetMeHandler(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig {
	event := ev.InnerEvent.Data.(*slackevents.AppMentionEvent)

	for index, ch := range botConfig.Channels {
		if ch.ChannelId == event.Channel {
			requiredChannel := &botConfig.Channels[index]

			requiredChannel.RecentUsers = requiredChannel.RecentUsers[:0]

			_, _, _ = client.PostMessage(event.Channel, slack.MsgOptionText("https://media.giphy.com/media/LOoaJ2lbqmduxOaZpS/giphy.gif", false))
			_, _, _ = client.PostMessage(event.Channel, slack.MsgOptionText("All recent users reset!", false))

		}
	}

	botConfig.SaveCurrentState(util.OsFileUtility{})
	return botConfig
}
