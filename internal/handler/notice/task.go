package notice

import (
	"log"

	"github.com/gin-gonic/gin"
)

var HubIns = NewHub()

// serveWs handles websocket requests from the peer.
func ServeWs(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: HubIns, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	client.readPump()
}
