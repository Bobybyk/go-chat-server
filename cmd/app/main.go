package main

import (
	// gestion console
	"fmt"
	// gestion des logs
	"log"
	// gestion des requêtes HTTP
	"net/http"
	// gestion des mutex (concurence)
	"sync"
	// gestion des websockets
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
var clients = make(map[*websocket.Conn]bool)
// Gère les messages à diffuser
var broadcast = make(chan Message)
// Gère l'accès concurrent aux clients
var mutex sync.Mutex

// Structure du message échangé
type Message struct {
	Content string `json:"content"`
}

/** Configuration du gestionnaire d'itinéraire pour la connexion WebSocket, 
 * démarrage d'une goroutine pour gérer les messages entrants, 
 * sert les fichiers statiques ou l'interface utilisateur (plus tard), 
 * écoute sur le port 8081
 */
func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	// Serve static files or use your own frontend
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server started on :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/**
 * Met à jour la carte des clients, 
 * lit les messages JSON des clients et les transmet à la goroutine de diffusion,
 * et gère les erreurs de connexion
 */
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		broadcast <- msg
	}
}

/**
 * Diffuse les messages aux clients via la carte des clients,
 * boucle en continu pour recevoir les messages de la goroutine handleConnections
 */
func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
