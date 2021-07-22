package azuredevops

import (
	"fmt"
	"meeseeks/entity"
	"regexp"

	"github.com/slack-go/slack/slackevents"
)

func RetrieveTickets(meeseeks *entity.MeeseeksSlack, message *slackevents.MessageEvent) {
	ticketDetect, _ := regexp.Compile(`#\d{4}`)
	possibleTickets := ticketDetect.FindAllString(message.Text, -1)

	msg := fmt.Sprintf("\n\n\n\n\nI saw a message from %q in %q that looked like it contained a ticket number:\n", message.User, message.Channel)

	if len(possibleTickets) > 0 {
		for _, tkt := range possibleTickets {
			msg += fmt.Sprintf("Ticket: %s \n", tkt)
		}
	}

	fmt.Printf(msg)
}
