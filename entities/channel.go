package entities

type Channel struct {
	ChannelName   string
	ChannelId     string
	EligibleUsers []string
	ExemptUsers   []string
	RecentUsers   []string
}

func (channel Channel)