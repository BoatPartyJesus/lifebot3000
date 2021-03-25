package main

import (
	"awesomeProject/configuration"
	"awesomeProject/entities"
	"awesomeProject/handlers"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	fmt.Println("Starting a thing")

	//configFile := whatever arg 1 is...
	//botConfig := configuration.LoadConfiguration(configFile)
	botConfig := entities.LifeBotConfig{
		AppToken: "",
		BotToken: "",
		Channels: nil,
	}

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
						if strings.Contains(ev.Text, "addme") == true {
							command := "addme"

							handler := handlers[command]
							if handler == nil {
								return
							}

							handler(eventsAPIEvent, api)
						}
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
