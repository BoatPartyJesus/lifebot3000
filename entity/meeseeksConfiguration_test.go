package entity

import (
	"encoding/json"
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fileMock struct {
	mock.Mock
}

func (f *fileMock) ReadFile(filepath string) ([]byte, error) {
	args := f.Called(filepath)
	return args.Get(0).([]byte), args.Error(1)
}

func (f *fileMock) WriteFile(filepath string, fileContents []byte) error {
	_ = f.Called(filepath)
	return nil
}

func (f *fileMock) CreateFile(filepath string) error {
	_ = f.Called(filepath)
	return nil
}

func (f *fileMock) DoesFileExist(filepath string) bool {
	args := f.Called(filepath)
	return args.Get(0).(bool)
}

func TestRetrieveOrCreateChannel_NewChannel(t *testing.T) {
	type args struct {
		foundChannel  slack.Channel
		knownChannels []Channel
	}
	type test struct {
		name string
		args args
		want []Channel
	}

	shouldAddNew := test{
		name: "ShouldAddNewChannel",
		args: args{
			foundChannel: slack.Channel{
				GroupConversation: slack.GroupConversation{Conversation: slack.Conversation{ID: "newId", NameNormalized: "newChannel"}},
			},
			knownChannels: []Channel{{
				ChannelName: "oldChannel",
				ChannelId:   "oldId",
			}},
		},
		want: []Channel{{
			ChannelName: "oldChannel",
			ChannelId:   "oldId",
		}, {
			ChannelName: "newChannel",
			ChannelId:   "newId",
		}},
	}

	t.Run(shouldAddNew.name, func(t *testing.T) {
		got := RetrieveOrCreateChannel(shouldAddNew.args.foundChannel, shouldAddNew.args.knownChannels)

		assert.Equal(t, len(shouldAddNew.want), len(got), "array lengths should be the same")

		for i, ch := range got {
			assert.Equal(t, shouldAddNew.want[i].ChannelId, ch.ChannelId, "channel ID should match")
			assert.Equal(t, shouldAddNew.want[i].ChannelName, ch.ChannelName, "channel name should match")
		}
	})
}

func TestRetrieveOrCreateChannel_ExistingChannel(t *testing.T) {
	type args struct {
		foundChannel  slack.Channel
		knownChannels []Channel
	}
	type test struct {
		name string
		args args
		want []Channel
	}

	shouldNotAddNew := test{
		name: "ShouldNotAddNewChannel",
		args: args{
			foundChannel: slack.Channel{
				GroupConversation: slack.GroupConversation{Conversation: slack.Conversation{ID: "oldId", NameNormalized: "oldChannel"}},
			},
			knownChannels: []Channel{{
				ChannelName: "oldChannel",
				ChannelId:   "oldId",
			}},
		},
		want: []Channel{{
			ChannelName: "oldChannel",
			ChannelId:   "oldId",
		},
		},
	}

	t.Run(shouldNotAddNew.name, func(t *testing.T) {
		got := RetrieveOrCreateChannel(shouldNotAddNew.args.foundChannel, shouldNotAddNew.args.knownChannels)

		assert.Equal(t, len(shouldNotAddNew.want), len(got), "array lengths should be the same")

		for i, ch := range got {
			assert.Equal(t, shouldNotAddNew.want[i].ChannelId, ch.ChannelId, "channel ID should match")
			assert.Equal(t, shouldNotAddNew.want[i].ChannelName, ch.ChannelName, "channel name should match")
		}
	})
}

func TestMeeseeksConfig_SaveCurrentState_FileExists(t *testing.T) {

	fu := &fileMock{}

	fu.On("DoesFileExist", "godThisGameIsCrap").Return(true)
	fu.On("WriteFile", "godThisGameIsCrap").Return(nil)

	meeseeks := MeeseeksConfig{FileLocation: "godThisGameIsCrap"}
	t.Run("ShouldMakeExpectedFileCalls", func(t *testing.T) {
		meeseeks.SaveCurrentState(fu)
	})

	fu.AssertCalled(t, "WriteFile", "godThisGameIsCrap")
}

func TestMeeseeksConfig_SaveCurrentState_FileDoesNotExist(t *testing.T) {

	fu := &fileMock{}

	fu.On("DoesFileExist", "godThisGameIsCrap").Return(false)
	fu.On("CreateFile", "godThisGameIsCrap").Return(nil)
	fu.On("WriteFile", "godThisGameIsCrap").Return(nil)

	meeseeks := MeeseeksConfig{FileLocation: "godThisGameIsCrap"}

	meeseeks.SaveCurrentState(fu)

	fu.AssertCalled(t, "CreateFile", "godThisGameIsCrap")
	fu.AssertCalled(t, "WriteFile", "godThisGameIsCrap")
}

func TestMeeseeksConfig_LoadState(t *testing.T) {

	fu := &fileMock{}

	storedState, _ := json.MarshalIndent(MeeseeksConfig{}, "", "  ")

	fu.On("DoesFileExist", "pleaseGodMakeItStop").Return(true)
	fu.On("ReadFile", "pleaseGodMakeItStop").Return(storedState, nil)

	actual := LoadConfigurationStateFromFile(fu, "pleaseGodMakeItStop")

	assert.Equal(t, "pleaseGodMakeItStop", actual.FileLocation, "File location should be set correctly for older file")
}
