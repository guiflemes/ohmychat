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

type OhMyChatConfig struct {
	Api    Api
	Worker Worker
}
