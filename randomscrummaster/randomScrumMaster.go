package randomscrummaster

import (
	"fmt"
	"meeseeks/entity"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/slack-go/slack"
)

func PickWinnerByChannel(meeseeks *entity.MeeseeksSlack, config *entity.MeeseeksConfig, requiredChannel string) {
	for index, channel := range config.Channels {
		if channel.ChannelId == requiredChannel {

			luckyWinner, _ := channel.PickAWinner(meeseeks)

			if luckyWinner == "" {
				fmt.Println("No eligible users to pick.")

				headerText := slack.NewTextBlockObject("mrkdwn", "There's nobody to pick :cry:", false, false)
				header := slack.NewSectionBlock(headerText, nil, nil)

				volunteerBtnTxt := slack.NewTextBlockObject("plain_text", "I'll do it!", false, false)
				volunteerBtn := slack.NewButtonBlockElement("rsm_volunteer", "rsm_volunteer", volunteerBtnTxt)

				actionBlock := slack.NewActionBlock("", volunteerBtn)

				msg := slack.MsgOptionBlocks(header, actionBlock)

				_, _, _ = meeseeks.Slack.PostMessage(channel.ChannelId, slack.MsgOptionText("https://media.giphy.com/media/d2lcHJTG5Tscg/giphy.gif", false))
				_, _, _ = meeseeks.Slack.PostMessage(channel.ChannelId, slack.MsgOptionText("There's nobody to pick :(", false), msg)

			} else {
				fmt.Println("Winner:" + luckyWinner)
				config.Channels[index].RecentUsers = channel.AddRecentUser(luckyWinner)
				config.SaveCurrentState()

				winnerMessage := fmt.Sprintf("Today's scrum master is <@%s>! :tada:\n", luckyWinner)

				headerText := slack.NewTextBlockObject("mrkdwn", winnerMessage, false, false)
				header := slack.NewSectionBlock(headerText, nil, nil)

				rerollBtnTxt := slack.NewTextBlockObject("plain_text", "Re-Roll", false, false)
				rerollBtn := slack.NewButtonBlockElement("rsm_reroll", "rsm_reroll", rerollBtnTxt)

				volunteerBtnTxt := slack.NewTextBlockObject("plain_text", "I'll do it!", false, false)
				volunteerBtn := slack.NewButtonBlockElement("rsm_volunteer", "rsm_volunteer", volunteerBtnTxt)

				actionBlock := slack.NewActionBlock("", rerollBtn, volunteerBtn)

				msg := slack.MsgOptionBlocks(header, actionBlock)

				_, _, _ = meeseeks.Slack.PostMessage(channel.ChannelId, slack.MsgOptionText(winnerMessage, false), msg)
			}
		} else {
			continue
		}
	}
}

var pickWinners = func(meeseeks *entity.MeeseeksSlack, config *entity.MeeseeksConfig) {
	fmt.Println("Picking a random scrum master...")

	channels := config.Channels

	for _, channel := range channels {
		PickWinnerByChannel(meeseeks, config, channel.ChannelId)
	}
}

func PickAWinnerCron(meeseeks *entity.MeeseeksSlack, config *entity.MeeseeksConfig) {
	s := gocron.NewScheduler(time.UTC)
	// s.Cron("0 9 * * 1-5").Do(pickWinners(config)) // 9am daily
	s.Cron("* * * * *").Do(pickWinners, meeseeks, config) // every min
	s.StartBlocking()
}
