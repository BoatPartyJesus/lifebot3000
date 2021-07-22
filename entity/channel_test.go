package entity

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type meeseeksMock struct {
	mock.Mock
}

func (m *meeseeksMock) GetAvailableUsers(users []string) []string {
	args := m.Called(users)
	return args.Get(0).([]string)
}

func (m *meeseeksMock) New(botConfig *MeeseeksConfig) {

}

func TestChannel_IsEligibleUser(t *testing.T) {
	type args struct {
		x string
	}
	tests := []struct {
		name    string
		channel Channel
		args    args
		want    bool
	}{
		{name: "ShouldBeEligible", args: args{x: "yes, dear boy"}, channel: Channel{EligibleUsers: []string{"yes, dear boy"}}, want: true},
		{name: "ShouldNotBeEligible", args: args{x: "certainly not"}, channel: Channel{EligibleUsers: []string{"yes, dear boy"}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.channel.IsEligibleUser(tt.args.x); got != tt.want {
				assert.Equal(t, tt.want, got, tt.name)
			}
		})
	}
}

func TestChannel_AddEligibleUser(t *testing.T) {
	type args struct {
		x string
	}
	tests := []struct {
		name    string
		channel Channel
		args    args
		want    []string
	}{
		{name: "ShouldAddEligibleUser", args: args{x: "Off you pop!"}, channel: Channel{EligibleUsers: []string{"yes, dear boy"}}, want: []string{"yes, dear boy", "Off you pop!"}},
		{name: "ShouldNotAddDuplicateUser", args: args{x: "Off you pop!"}, channel: Channel{EligibleUsers: []string{"Off you pop!"}}, want: []string{"Off you pop!"}},
		{name: "ShouldAllowOverThreeUsers", args: args{x: "u5"}, channel: Channel{EligibleUsers: []string{"u1", "u2", "u3", "u4"}}, want: []string{"u1", "u2", "u3", "u4", "u5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.channel.AddEligibleUser(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got, tt.name)
			}
		})
	}
}

func TestChannel_RemoveEligibleUser(t *testing.T) {
	type args struct {
		x string
	}
	test := struct {
		name         string
		channel      Channel
		args         args
		wantEligible []string
		wantRecent   []string
	}{
		name:         "ShouldRemoveEligibleUser",
		args:         args{x: "Off you pop!"},
		channel:      Channel{EligibleUsers: []string{"yes, dear boy", "Off you pop!"}, RecentUsers: []string{"Off you pop!"}},
		wantEligible: []string{"yes, dear boy"},
		wantRecent:   []string{},
	}

	t.Run(test.name, func(t *testing.T) {
		gotEligible, gotRecent := test.channel.RemoveEligibleUser(test.args.x)

		if !assert.ObjectsAreEqual(test.wantEligible, gotEligible) {
			t.Errorf("Channel.RemoveEligibleUser() gotEligible = %v, wantEligible %v", gotEligible, test.wantEligible)
		}
		if !assert.ObjectsAreEqual(gotRecent, test.wantRecent) {
			t.Errorf("Channel.RemoveEligibleUser() gotRecent = %v, wantRecent %v", gotRecent, test.wantRecent)
		}
	})
}

func TestChannel_AddRecentUser(t *testing.T) {
	type args struct {
		x string
	}
	tests := []struct {
		name    string
		channel Channel
		args    args
		want    []string
	}{
		{name: "ShouldAddRecentUser", args: args{x: "Off you pop!"}, channel: Channel{RecentUsers: []string{"yes, dear boy"}}, want: []string{"yes, dear boy", "Off you pop!"}},
		{name: "ShouldNotAddDuplicateUser", args: args{x: "Off you pop!"}, channel: Channel{RecentUsers: []string{"Off you pop!"}}, want: []string{"Off you pop!"}},
		{name: "ShouldRotateMoreThanThreeUsers", args: args{x: "u4"}, channel: Channel{RecentUsers: []string{"u1", "u2", "u3"}}, want: []string{"u2", "u3", "u4"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.channel.AddRecentUser(tt.args.x)
			if !assert.ObjectsAreEqual(tt.want, got) {
				t.Errorf("Channel.AddRecentUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_remove(t *testing.T) {
	type args struct {
		s []string
		x string
	}
	test := struct {
		name string
		args args
		want []string
	}{
		name: "ShouldRemoveElement", args: args{s: []string{"a", "b", "c"}, x: "b"}, want: []string{"a", "c"},
	}
	t.Run(test.name, func(t *testing.T) {
		got := remove(test.args.s, test.args.x)
		if !assert.ObjectsAreEqual(test.want, got) {
			t.Errorf("remove() = %v, want %v", got, test.want)
		}
	})
}

func Test_find(t *testing.T) {
	type args struct {
		s []string
		x string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "ShouldFindElement", args: args{s: []string{"a", "b", "c"}, x: "b"}, want: true},
		{name: "ShouldFindElement", args: args{s: []string{"a", "b", "c"}, x: "d"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, find(tt.args.s, tt.args.x), tt.name)
		})
	}
}

func TestChannel_PickAWinner(t *testing.T) {

	msks := new(meeseeksMock)

	msks.On("GetAvailableUsers", []string{"a", "b", "c"}).Return([]string{"a"})
	msks.On("GetAvailableUsers", []string{"b", "c"}).Return([]string{})

	type args struct {
		meeseeks IMeeseeksSlack
	}
	tests := []struct {
		name    string
		channel Channel
		args    args
		want    string
		wantErr bool
	}{
		{"ShouldPickAvailableUser", Channel{EligibleUsers: []string{"a", "b", "c"}}, args{meeseeks: msks}, "a", false},
		{"ShouldNotPickUnavailableUser", Channel{EligibleUsers: []string{"b", "c"}}, args{meeseeks: msks}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.channel.PickAWinner(tt.args.meeseeks)
			if tt.wantErr {
				assert.Equal(t, fmt.Errorf("no available users in channel"), err)
				return
			}

			assert.Equal(t, tt.want, got, tt.name)
		})
	}
}
