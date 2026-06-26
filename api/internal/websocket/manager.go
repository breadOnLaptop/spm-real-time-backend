package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for MVP
	},
}

type Manager struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		connections: make(map[*websocket.Conn]bool),
	}
}

func (m *Manager) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
		return
	}

	m.mu.Lock()
	m.connections[conn] = true
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.connections, conn)
		m.mu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (m *Manager) Broadcast(message []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for conn := range m.connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Broadcast error: %v", err)
			conn.Close()
			delete(m.connections, conn)
		}
	}
}
