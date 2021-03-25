package entities

type LifeBotConfig struct {
	AppToken string
	BotToken string
	Channels []Channel
}

type Channel struct {
	ChannelName   string
	EligibleUsers []string
	ExemptUsers   []string
	RecentUsers   []string
}
