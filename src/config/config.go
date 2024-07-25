package config

type Worker struct {
	Enabled bool
	Number  int
}

type Api struct {
	Enabled bool
	Port    int
}

const (
	MessageOmitted string = "omitted"
)

type Provider string

const (
	Telegram Provider = "telegram"
	Cli      Provider = "cli"
)

type Connector struct {
	Provider Provider
}

type ChatDatabase struct {
	Kind     string
	Settings map[string]any
}

type OhMyChatConfig struct {
	Api          Api
	Worker       Worker
	Connector    Connector
	ChatDatabase ChatDatabase
}
