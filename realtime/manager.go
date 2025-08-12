// wsmanager/manager.go
package wsmanager

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	Conn *websocket.Conn
}

type WSManager struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]bool
}

var manager = &WSManager{
	clients: make(map[*websocket.Conn]bool),
}

func AddClient(conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.clients[conn] = true
}

func RemoveClient(conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.clients, conn)
}

func Broadcast(message string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	for client := range manager.clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			client.Close()
			delete(manager.clients, client)
		}
	}
}
