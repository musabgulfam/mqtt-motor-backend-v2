package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	wsmanager "github.com/musabgulfam/pumplink-backend/realtime"
	"github.com/musabgulfam/pumplink-backend/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Expect the first message to be the JWT token
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read auth token:", err)
		return
	}
	token := string(msg)

	_, err = utils.ValidateJWT(token)
	if err != nil {
		log.Println("Invalid token, closing connection")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Invalid token"))
		return
	}

	wsmanager.AddClient(conn)          // Add the new client to the manager
	defer wsmanager.RemoveClient(conn) // Ensure client is removed on disconnect
	fmt.Println("Client connected:", conn.RemoteAddr())

	for {
		_, _, err := conn.ReadMessage() // Keep connection alive
		if err != nil {
			wsmanager.RemoveClient(conn) // Remove client on error
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			break
		}
	}
}
