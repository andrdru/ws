package ws

import (
	"encoding/json"
)

type (
	Message interface {
		json.Marshaler
		json.Unmarshaler

		GetType() MessageType
		IsSenderBroadcastSkipped() bool
	}

	Client interface {
		GetHub() Hub
		GetID() string
		SetID(clientID string)
		GetMessageChan() chan Message
	}

	Hub interface {
		Run()
		Register(c Client)
		Unregister(c Client)
		SendBroadcast(m HubMessage)
		SendPersonal(m HubMessage)
		ClientReadHandler() ClientReadHandler
		SetClientReadHeader(handler ClientReadHandler)
	}
)
