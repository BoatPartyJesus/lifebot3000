package main

import (
	"flag"
	"fmt"
	"log"
	"meeseeks/configuration"
	"meeseeks/entities"
	"meeseeks/handlers"
	channelHelper "meeseeks/helpers"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func pickAWinnerCron(api *slack.Client, config *entities.LifeBotConfig) {
	s := gocron.NewScheduler(time.UTC)
	// s.Cron("0 9 * * 1-5").Do(pickWinners(config)) // 9am daily
	s.Cron("* * * * *").Do(pickWinners, api, config) // every min
	s.StartBlocking()
}

var pickWinners = func(api *slack.Client, config *entities.LifeBotConfig) {
	fmt.Println("Picking a random scrum master...")

	channels := config.Channels

	for index, channel := range channels {

		var luckyWinner string

		for {
			user, _ := channelHelper.PseudoRandomSelect(channel.EligibleUsers, channel.RecentUsers)
			userProfile, _ := api.GetUserProfile(user, true)

			if !channelHelper.Find(config.ExcludedSlackStatuses, userProfile.StatusText) {
				luckyWinner = user
				break
			}
		}

		if luckyWinner == "" {
			fmt.Println("No eligible users to pick.")
		} else {
			fmt.Println("Winner:" + luckyWinner)
			config.Channels[index].RecentUsers = channelHelper.AddRecentUser(channel.RecentUsers, luckyWinner)
			config.SaveCurrentState()
		}
	}
}

func main() {
	fmt.Println("Starting a thing")

	var configFile string
	flag.StringVar(&configFile, "c", "config", "Specify config location")
	flag.Parse()

	botConfig := configuration.LoadConfiguration(configFile)
	meeseeks := new(channelHelper.MeeseeksSlack)
	meeseeks.New(&botConfig)

	api := meeseeks.Slack
	client := socketmode.New(
		meeseeks.Slack,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)))

	channelHandlers := map[string]func(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entities.LifeBotConfig) entities.LifeBotConfig{
		"addme":    handlers.ChannelAddMeHandler,
		"list":     handlers.ListHandler,
		"removeme": handlers.RemoveMeHandler,
	}
	// messageHandlers := map[string]func(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entities.LifeBotConfig) entities.LifeBotConfig{
	// 	"addme":    handlers.MessageAddMeHandler,
	// 	"list":     handlers.ListHandler,
	// 	"removeme": handlers.RemoveMeHandler,
	// }

	go pickAWinnerCron(api, &botConfig)

	go func() {
		for event := range client.Events {
			switch event.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Failed to connect to Slack")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to api")
			case socketmode.EventTypeHello:
				testResult, _ := api.AuthTest()
				fmt.Println(testResult)
				params := slack.GetConversationsForUserParameters{
					UserID:          testResult.UserID,
					Cursor:          "",
					Types:           []string{"private_channel", "public_channel"},
					Limit:           0,
					ExcludeArchived: false,
				}

				result, _, err := api.GetConversationsForUser(&params)
				if err != nil {
					fmt.Println(err)
				}

				for _, ch := range result {
					botConfig.Channels = channelHelper.RetrieveOrCreate(ch, botConfig.Channels)
				}

			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", event)
					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*event.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						idPattern, _ := regexp.Compile(`<@\w{11}>\W*`)
						args := strings.Split(idPattern.ReplaceAllLiteralString(ev.Text, ""), " ")

						if args[0] == "" {
							_, _, _ = api.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("<@%s>...\n", ev.User), false))
						}

						handler := channelHandlers[args[0]]
						if handler == nil {
							return
						}

						botConfig = handler(eventsAPIEvent, api, botConfig)
					// case *slackevents.MessageEvent: // Different interface, can't bundle nicely... yet...
					// 	idPattern, _ := regexp.Compile(`<@\w{11}>\W*`)
					// 	args := strings.Split(idPattern.ReplaceAllLiteralString(ev.Text, ""), " ")

					// 	if args[0] == "" {
					// 		_, _, _ = api.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("<@%s>...\n", ev.User), false))
					// 	}

					// 	handler := messageHandlers[args[0]]
					// 	if handler == nil {
					// 		return
					// 	}

					// 	botConfig = handler(eventsAPIEvent, api, botConfig)
					case *slackevents.MemberJoinedChannelEvent:
						fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
					}
				default:
					client.Debugf("unsupported Events API event received")
				}
			default:
				fmt.Printf("Unexpected event type received: %s\n\n", event.Type)
			}
		}
	}()

	client.Run()

	// signalChannel := make(chan os.Signal)
	// signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	// <-signalChannel

	//configuration.SaveConfiguration(configFile, botConfig)

}
