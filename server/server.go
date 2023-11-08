package server

import (
	"fmt"
	// Importation du gestionnaire de temps
	"time"
	// Importation mutex
	"sync"
	// Importation du gestionnaire de route
	"net/http"
	// Importation du gestionnaire websocket
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

// Structure d'un client
type Client struct {
	Conn      *websocket.Conn
	Username  string
}

// types de messages
const (
	MessageTypeNormal    = "normal"
	MessageTypeJoin      = "join"
	MessageTypeLeave     = "leave"
)

// Envoie tous les messages de l'historique à la connexion
func sendChatHistory(conn *websocket.Conn) {
    mutex.Lock()
    for _, message := range chatHistory {
        err := conn.WriteJSON(message)
        if err != nil {
            fmt.Printf("Error writing message: %v", err)
        }
    }
    mutex.Unlock()
}

// contient les messages déjà reçus
var chatHistory []Message