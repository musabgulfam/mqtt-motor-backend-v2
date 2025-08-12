package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

// Keep track of connected clients
var wsClients = make(map[*websocket.Conn]bool)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	wsClients[conn] = true
	fmt.Println("Client connected:", conn.RemoteAddr())

	for {
		_, _, err := conn.ReadMessage() // Keep connection alive
		if err != nil {
			delete(wsClients, conn)
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			break
		}
	}
}

// Send data to all connected WS clients
func BroadcastToWSClients(message string) {
	for client := range wsClients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			client.Close()
			delete(wsClients, client)
		}
	}
}
