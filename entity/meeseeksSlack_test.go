package entity

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeeseeksSlack_New(t *testing.T) {
	type args struct {
		botConfig *MeeseeksConfig
	}

	botConfig := MeeseeksConfig{FileLocation: "flibble"}

	tests := []struct {
		name     string
		meeseeks *MeeseeksSlack
		args     args
		want     string
	}{
		{name: "ShouldSetExpectedConfig", meeseeks: &MeeseeksSlack{}, args: args{botConfig: &botConfig}, want: "flibble"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.meeseeks.New(tt.args.botConfig)
			assert.Equal(t, tt.want, tt.meeseeks.BotConfig.FileLocation, "correct config should be used")
		})
	}
}

func TestMeeseeksSlack_GetAvailableUsers(t *testing.T) {
	// Need to Mock for this - read up on GoMock
	type args struct {
		users []string
	}
	tests := []struct {
		name     string
		meeseeks *MeeseeksSlack
		args     args
		want     []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.meeseeks.GetAvailableUsers(tt.args.users); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MeeseeksSlack.GetAvailableUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arraySeek(t *testing.T) {
	type args struct {
		s []string
		x string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ShouldFindInArray", args{[]string{"a", "b", "c"}, "b"}, true},
		{"ShouldNotFindInArray", args{[]string{"a", "b", "c"}, "d"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, arraySeek(tt.args.s, tt.args.x), tt.name)
		})
	}
}
