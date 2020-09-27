package ws

import (
	"github.com/gofrs/uuid"
)

var _ Hub = &hub{}

func NewHub() *hub {
	return &hub{
		broadcast:  make(chan HubMessage, 1),
		personal:   make(chan HubMessage, 1),
		register:   make(chan Client),
		unregister: make(chan Client),
		clients:    make(map[Client]struct{}),
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			client.SetID(uuid.Must(uuid.NewV4()).String())
			h.clients[client] = struct{}{}

			if h.RegisterRespHandler != nil {
				client.GetMessageChan() <- h.RegisterRespHandler(client)
			}

			if h.NewOnlineHandler != nil {
				h.broadcast <- HubMessage{
					Message: h.NewOnlineHandler(client),
					From:    client,
				}
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.GetMessageChan())

				if h.OfflineHandler != nil {
					h.broadcast <- HubMessage{
						Message: h.OfflineHandler(client),
						From:    client,
					}
				}
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				if message.Message.IsSenderBroadcastSkipped() && message.From == client {
					continue
				}

				select {
				case client.GetMessageChan() <- message.Message:
				default:
					close(client.GetMessageChan())
					delete(h.clients, client)
				}
			}
		case message := <-h.personal:
			var ok bool
			if _, ok = h.clients[message.To]; !ok {
				continue
			}

			select {
			case message.To.GetMessageChan() <- message.Message:
			default:
				close(message.To.GetMessageChan())
				delete(h.clients, message.To)
			}
		}
	}
}

func (h *hub) SendBroadcast(m HubMessage) {
	h.broadcast <- m
}

func (h *hub) SendPersonal(m HubMessage) {
	h.personal <- m
}

func (h *hub) Unregister(c Client) {
	h.unregister <- c
}

func (h *hub) Register(c Client) {
	h.register <- c
}

func (h *hub) ClientReadHandler() ClientReadHandler {
	return h.clientReadHandler
}

func (h *hub) SetClientReadHeader(handler ClientReadHandler) {
	h.clientReadHandler = handler
}
