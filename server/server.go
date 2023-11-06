package server

import (
	"time"
	"sync"
	"net/http"
	"github.com/gorilla/websocket"
)

// Configuration des websockets
var upgrader = websocket.Upgrader{
	// Autorise toutes les connexions (à modifier pour ajouter des restrictions)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Gère les clients connectés
var clients = make(map[*websocket.Conn]string)
// Gère les messages à diffuser
var broadcast = make(chan Message)
// Gère l'accès concurrent aux clients
var mutex sync.Mutex

// Structure du message échangé
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Type     string `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

type Client struct {
	Conn      *websocket.Conn
	Username  string
}

const (
	MessageTypeNormal    = "normal"
	MessageTypeJoin      = "join"
	MessageTypeLeave     = "leave"
)
