package ws

import (
	"github.com/gorilla/websocket"
)

//go:generate easyjson

type (
	MessageType       int
	EventHandler      func(c Client) Message
	MessageOption     func(*MessageOptions)
	ClientReadHandler func(client Client, bytes []byte) (err error)

	client struct {
		conn     *websocket.Conn
		clientID string
		send     chan Message

		Hub Hub

		ReadHandler ClientReadHandler
	}

	HubMessage struct {
		Message Message
		To      Client
		From    Client
	}

	hub struct {
		clients map[Client]struct{}

		broadcast chan HubMessage
		personal  chan HubMessage

		register   chan Client
		unregister chan Client

		RegisterRespHandler EventHandler
		NewOnlineHandler    EventHandler
		OfflineHandler      EventHandler

		clientReadHandler ClientReadHandler
	}

	MessageOptions struct {
		isSenderBroadcastSkipped bool
	}

	//easyjson:json
	BaseMessage struct {
		MessageType              MessageType `json:"type"`
		isSenderBroadcastSkipped bool
	}
)
