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
	Cli      MessageConnector = "cli"
)

type ResponseType int

const (
	OptionResponse ResponseType = iota
	TextResponse
)

type Meta struct {
	Data map[string]string `json:"data"`
}

func (m *Meta) Add(name, value string) {
	m.Data["name"] = value
}

func (m *Meta) Get(name string) string {
	value, ok := m.Data[name]
	if !ok {
		return ""
	}
	return value
}

type Option struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID           string           `json:"id"`
	Type         MessageType      `json:"type"`
	Service      MessageService   `json:"service"`
	Connector    MessageConnector `json:"connector"`
	ConnectorID  string           `json:"connector_id"`
	BotName      string           `json:"bot_name"`
	BotID        string           `json:"bot_id"`
	ChannelID    string           `json:"channel_id"`
	ChannelName  string           `json:"channel_name"`
	Input        string           `json:"input"`
	Output       string           `json:"output"`
	Error        string           `json:"error"`
	Options      []Option         `json:"options"`
	StartTime    int64            `json:"start_time"`
	EndTime      int64            `json:"end_time"`
	ResponseType ResponseType     `json:"response_type"`
	Meta         *Meta            `json:"meta"`
}

func NewMessage() Message {
	return Message{
		ID:        uuid.NewString(),
		StartTime: time.Now().Unix(),
		Meta:      &Meta{Data: make(map[string]string)},
	}
}
