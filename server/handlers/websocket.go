package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"go-grpc-chat/utils"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Username string
	Send     chan []byte
}

var (
	clients      = make(map[string]*Client) // 🔐 key = username
	clientsMutex sync.RWMutex
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type IncomingMessage struct {
	Type    string `json:"type"`    // "public" or "private"
	From    string `json:"from"`    // sender username
	To      string `json:"to"`      // recipient username (optional)
	Content string `json:"content"` // actual message
}

// ✅ Entry point for WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// ✅ Validate JWT token
	claims, err := utils.ParseJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	username := claims["username"].(string)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	client := &Client{
		Conn:     conn,
		Username: username,
		Send:     make(chan []byte),
	}

	// 🧠 Add client to the map
	clientsMutex.Lock()
	clients[username] = client
	clientsMutex.Unlock()

	go writeToClient(client)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg IncomingMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			fmt.Println("Invalid JSON:", err)
			continue
		}

		switch msg.Type {
		case "private":
			sendPrivateMessage(msg.From, msg.To, msg.Content)
		default:
			broadcastPublic(fmt.Sprintf("%s: %s", msg.From, msg.Content))
		}
	}

	// Clean up on disconnect
	clientsMutex.Lock()
	delete(clients, username)
	clientsMutex.Unlock()
}

func writeToClient(c *Client) {
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func sendPrivateMessage(from, to, content string) {
	clientsMutex.RLock()
	receiver, ok := clients[to]
	clientsMutex.RUnlock()

	if ok {
		message := fmt.Sprintf("🔒 [Private] %s: %s", from, content)
		receiver.Send <- []byte(message)
	}
}

func broadcastPublic(message string) {
	clientsMutex.RLock()
	for _, client := range clients {
		client.Send <- []byte(message)
	}
	clientsMutex.RUnlock()
}
