package entities

import "time"

type LifeBotConfig struct {
	AppToken    string
	BotToken    string
	Channels    []Channel
	LastUpdated time.Time
}

type Channel struct {
	ChannelName   string
	ChannelId     string
	EligibleUsers []string
	ExemptUsers   []string
	RecentUsers   []string
}
