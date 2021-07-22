package handler

import (
	"meeseeks/entity"
	"reflect"
	"testing"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func TestAddMeHandler(t *testing.T) {
	type args struct {
		ev        slackevents.EventsAPIEvent
		client    *slack.Client
		botConfig entity.MeeseeksConfig
		namedUser string
	}
	tests := []struct {
		name string
		args args
		want entity.MeeseeksConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddMeHandler(tt.args.ev, tt.args.client, tt.args.botConfig, tt.args.namedUser); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddMeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
