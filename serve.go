package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// serveWs handles websocket requests from the peer.
func ServeWs(hub Hub, upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(hub, conn)
	client.ReadHandler = hub.ClientReadHandler()
	client.Hub.Register(client)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
