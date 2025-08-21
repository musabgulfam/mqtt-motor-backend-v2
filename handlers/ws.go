package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	wsmanager "github.com/musabgulfam/pumplink-backend/services"
	"github.com/musabgulfam/pumplink-backend/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // Upgrade HTTP connection to WebSocket
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer conn.Close()

	// 1. Read the first message (expecting JWT)
	_, jwtMsg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Failed to read JWT from client: %v\n", err)
		return
	}

	// 2. Validate JWT (implement your own validation logic)
	_, err = utils.ValidateJWT(string(jwtMsg))
	if err != nil {
		log.Printf("Invalid JWT: %v\n", err)
		conn.WriteMessage(websocket.TextMessage, []byte("unauthorized"))
		return
	}

	// 3. Add client to manager and send latest status
	wsmanager.AddClient(conn)
	defer wsmanager.RemoveClient(conn) // Ensure client is removed on disconnect
	log.Printf("[WEBSOCKET] Client connected: %v\n", conn.RemoteAddr())

	// 4. Now enter your main read loop (if needed)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			wsmanager.RemoveClient(conn)
			log.Printf("[WEBSOCKET] read error: %v\n", err)
			break
		}
	}
	// wsmanager.Broadcast(fmt.Sprintf("New client connected: %v", conn.RemoteAddr().String())) // Notify all clients
}
