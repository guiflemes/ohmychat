package models

type MessageType int

const (
	MsgTypeUnknown MessageType = iota
	MsgTypeChannel
)

type MessageService int

const (
	MsgServiceUnknown MessageService = iota
	MsgServiceChat
)

type MessageRemote string

const (
	Telegram MessageRemote = "telegram"
)

type Message struct {
	ID          string
	Type        MessageType
	Service     MessageService
	Remote      MessageRemote
	ChannelID   string
	ChannelName string
	Input       string
	Output      string
	Error       string
}
