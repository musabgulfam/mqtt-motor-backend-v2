package services

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WSManager struct {
	clients      map[*websocket.Conn]bool
	mu           sync.Mutex
	latestStatus []byte
}

var manager = &WSManager{
	clients: make(map[*websocket.Conn]bool),
}

// AddClient registers a new client and sends the latest status if available
func AddClient(conn *websocket.Conn) {
	manager.mu.Lock()
	manager.clients[conn] = true
	latest := manager.latestStatus
	manager.mu.Unlock()

	// Send the latest status to the new client
	if latest != nil {
		if err := conn.WriteMessage(websocket.TextMessage, latest); err != nil {
			log.Printf("[WEBSOCKET] Failed to send latest status to new client: %v", err)
		}
	}
}

// RemoveClient unregisters a client
func RemoveClient(conn *websocket.Conn) {
	manager.mu.Lock()
	delete(manager.clients, conn)
	manager.mu.Unlock()
}

// Broadcast sends a message to all connected clients and stores it as the latest status
func Broadcast(msg string) {
	manager.mu.Lock()
	manager.latestStatus = []byte(msg)
	for client := range manager.clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Printf("[WEBSOCKET] broadcast error: %v", err)
			client.Close()
			delete(manager.clients, client)
		}
	}
	manager.mu.Unlock()
}
