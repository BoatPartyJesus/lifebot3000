package main

import (
	"flag"
	"fmt"
	"log"
	ad "meeseeks/azuredevops"
	"meeseeks/entity"
	"meeseeks/handler"
	nlp "meeseeks/naturallanguageprocessing"
	rsm "meeseeks/randomscrummaster"
	"meeseeks/util"
	"os"
	"regexp"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	var botUserId string
	var configFile string
	flag.StringVar(&configFile, "c", "config", "Specify config location")
	flag.Parse()

	botConfig := entity.LoadConfigurationStateFromFile(util.OsFileUtility{}, configFile)

	meeseeks := new(entity.MeeseeksSlack)
	meeseeks.New(&botConfig)

	api := meeseeks.Slack

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)))

	//refactor to map[[]string]func
	// ["add", "include"] : handler.blah
	// or map aliases somewhere else...

	handlerAliases := map[string]func(ev slackevents.EventsAPIEvent, client *slack.Client, botConfig entity.MeeseeksConfig, namedUser string) entity.MeeseeksConfig{
		"add":    handler.AddMeHandler,
		"give":   handler.ListHandler,
		"remove": handler.RemoveMeHandler,
		"reset":  handler.ResetMeHandler,
	}

	// go rsm.PickAWinnerCron(meeseeks, &botConfig)

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
				botUserId = testResult.UserID

				params := slack.GetConversationsForUserParameters{
					UserID:          botUserId,
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

				// fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*event.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						event := innerEvent.Data.(*slackevents.AppMentionEvent)

						mentionChannel, _ := api.GetConversationInfo(ev.Channel, false)

						botConfig.Channels = entity.RetrieveOrCreateChannel(*mentionChannel, botConfig.Channels)

						idPattern, _ := regexp.Compile(fmt.Sprintf(`<@%s>\W*`, botUserId))
						mention := idPattern.ReplaceAllLiteralString(ev.Text, "")

						intent := nlp.IntentMap(mention)

						if intent.Action == "remove" && intent.Target == "recent" {
							meeseeks.Slack.PostEphemeral(event.Channel, event.User, slack.MsgOptionText("Yeah, you're not allowed to do that...", false))
							continue
						}

						handler := handlerAliases[intent.Action]
						if handler == nil {
							meeseeks.Slack.PostEphemeral(event.Channel, event.User, slack.MsgOptionText("Yeah, I dunno what the hell to do with that...", false))
							continue
						}
						botConfig = handler(eventsAPIEvent, api, botConfig, intent.Target)

					case *slackevents.MemberJoinedChannelEvent:
						// Prompt new joiner to be included in RSM, buttons for yes and no
						fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
						welcomeMessage := fmt.Sprintf("Hey <@%s>! I don't think I've seen you before! Should I include you in the random scrum master draw?\n", ev.User)

						headerText := slack.NewTextBlockObject("mrkdwn", welcomeMessage, false, false)
						header := slack.NewSectionBlock(headerText, nil, nil)

						acceptBtnTxt := slack.NewTextBlockObject("plain_text", "Yeah, sure!", false, false)
						acceptBtn := slack.NewButtonBlockElement("rsm_accept", "rsm_accept", acceptBtnTxt)

						declineBtnTxt := slack.NewTextBlockObject("plain_text", "Aww hell no!", false, false)
						declineBtn := slack.NewButtonBlockElement("rsm_decline", "rsm_decline", declineBtnTxt)

						actionBlock := slack.NewActionBlock("", acceptBtn, declineBtn)

						msg := slack.MsgOptionBlocks(header, actionBlock)

						_, _, _ = meeseeks.Slack.PostMessage(ev.Channel, slack.MsgOptionText(welcomeMessage, false), msg)

					case *slackevents.MessageEvent:

						ticketDetect, _ := regexp.Compile(`#\d{4}`)

						if ticketDetect.MatchString(ev.Text) {
							ad.RetrieveTickets(meeseeks, ev)
						}
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
					requiredAction := callback.ActionCallback.BlockActions[0].ActionID

					if requiredAction == "rsm_accept" {
						newUserChannel := callback.Container.ChannelID
						newUser := callback.User.ID

						for index, ch := range botConfig.Channels {
							if ch.ChannelId == newUserChannel {
								requiredChannel := &botConfig.Channels[index]

								var message string

								if ch.IsEligibleUser(newUser) {
									message = fmt.Sprintf("Huh... You are already in %s", requiredChannel.ChannelName)
								} else {
									requiredChannel.EligibleUsers = ch.AddEligibleUser(newUser)
									message = fmt.Sprintf("Cool, I'll add <@%s> to %s", newUser, requiredChannel.ChannelName)
								}

								_, _, _ = meeseeks.Slack.PostMessage(requiredChannel.ChannelId, slack.MsgOptionText(message, false))
							}
						}
					}

					if requiredAction == "rsm_decline" {
						newUserChannel := callback.Container.ChannelID

						message := "Fair enough! You can ask me to add you later if you change your mind!"
						_, _, _ = meeseeks.Slack.PostMessage(newUserChannel, slack.MsgOptionText(message, false))
					}

					if requiredAction == "rsm_reroll" {
						rerolledChannel := callback.Container.ChannelID
						rsm.PickWinnerByChannel(meeseeks, &botConfig, rerolledChannel, true)
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
}
