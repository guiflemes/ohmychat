package models

import (
	"time"

	"github.com/google/uuid"
)

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

type MessageConnector string

const (
	Telegram MessageConnector = "telegram"
	Test     MessageConnector = "testConn"
)

type ResponseType int

const (
	OptionResponse ResponseType = iota
)

type Meta struct {
	data map[string]string
}

func (m *Meta) Add(name, value string) {
	m.data["name"] = value
}

func (m *Meta) Get(name string) string {
	value, ok := m.data[name]
	if !ok {
		return ""
	}
	return value
}

type Option struct {
	ID   string
	Name string
}

type Message struct {
	ID           string
	Type         MessageType
	Service      MessageService
	Connector    MessageConnector
	ConnectorID  string
	BotName      string
	BotID        string
	ChannelID    string
	ChannelName  string
	Input        string
	Output       string
	Error        string
	Options      []Option
	StartTime    int64
	EndTime      int64
	ResponseType ResponseType
	Meta         *Meta
}

func NewMessage() Message {
	return Message{
		ID:        uuid.NewString(),
		StartTime: time.Now().Unix(),
		Meta:      &Meta{data: make(map[string]string)},
	}
}
