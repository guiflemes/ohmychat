package config

type Provider string

const (
	Telegram Provider = "telegram"
	Cli      Provider = "cli"
)

type Connector struct {
	Provider Provider
}

type OhMyChatConfig struct {
	Connector Connector
}
