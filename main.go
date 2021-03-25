package main

import (
	"awesomeProject/configuration"
	"awesomeProject/handlers"
	"flag"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

func main() {
	fmt.Println("Starting a thing")

	var configFile string
	flag.StringVar(&configFile, "c", "config", "Specify config location")
	flag.Parse()

	botConfig := configuration.LoadConfiguration(configFile)

	api := slack.New(
		botConfig.BotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(botConfig.AppToken))

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)))

	handlers := map[string]func(ev slackevents.EventsAPIEvent, client *slack.Client){
		"addme": handlers.AddMeHandler,
		"list":  handlers.ListHandler,
	}

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
				fmt.Println("Oh, HELLO THERE")

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
				fmt.Println(result)
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
						idPattern, _ := regexp.Compile("<@\\w{11}>\\W*")
						args := strings.Split(idPattern.ReplaceAllLiteralString(ev.Text, ""), " ")

						handler := handlers[args[0]]
						if handler == nil {
							return
						}

						handler(eventsAPIEvent, api)
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

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel

	//configuration.SaveConfiguration(configFile, botConfig)

}
