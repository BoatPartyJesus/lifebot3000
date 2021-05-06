package main

import (
	"flag"
	"fmt"
	"log"
	"meeseeks/configuration"
	"meeseeks/entity"
	"meeseeks/handler"
	rsm "meeseeks/randomscrummaster"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config", "Specify config location")
	flag.Parse()

	botConfig := configuration.LoadConfiguration(configFile)
	meeseeks := new(entity.MeeseeksSlack)
	meeseeks.New(&botConfig)

	api := meeseeks.Slack
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)))

	channelhandler := map[string]func(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig) entity.MeeseeksConfig{
		"addme":    handler.ChannelAddMeHandler,
		"list":     handler.ListHandler,
		"removeme": handler.RemoveMeHandler,
	}

	// messagehandler := map[string]func(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.LifeBotConfig) entity.LifeBotConfig{
	// 	"addme":    handler.MessageAddMeHandler,
	// 	"list":     handler.ListHandler,
	// 	"removeme": handler.RemoveMeHandler,
	// }

	go rsm.PickAWinnerCron(meeseeks, &botConfig)

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
					botConfig.Channels = entity.RetrieveOrCreateChannel(ch, botConfig.Channels)
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

						handler := channelhandler[args[0]]
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

					// 	handler := messagehandler[args[0]]
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
			case socketmode.EventTypeInteractive:
				callback, ok := event.Data.(slack.InteractionCallback)
				if !ok {
					fmt.Printf("Ignored %+v\n", event)

					continue
				}

				fmt.Printf("Interaction received: %+v\n", callback)

				var payload interface{}

				switch callback.Type {
				case slack.InteractionTypeBlockActions:
					// See https://api.slack.com/apis/connections/socket-implement#button

					requiredAction := callback.ActionCallback.BlockActions[0].ActionID

					if requiredAction == "rsm_reroll" {
						rerolledChannel := callback.Container.ChannelID
						rsm.PickWinnerByChannel(meeseeks, &botConfig, rerolledChannel)
					}

					if requiredAction == "rsm_volunteer" {
						volunteerChannel := callback.Container.ChannelID
						volunteer := callback.User.ID

						_, _, _ = meeseeks.Slack.PostMessage(volunteerChannel, slack.MsgOptionText("https://media.giphy.com/media/54JLdulN5BOwM/giphy.gif", false))
						_, _, _ = meeseeks.Slack.PostMessage(volunteerChannel, slack.MsgOptionText(fmt.Sprintf("<@%s> volunteered to be today's Random Scrum Master! \n", volunteer), false))
					}

				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
				case slack.InteractionTypeDialogSubmission:
				default:

				}

				client.Ack(*event.Request, payload)
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
