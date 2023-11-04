package server

import (
	"fmt"
	"log"
	"net/http"
)

/**
 * Met à jour la carte des clients, 
 * lit les messages JSON des clients et les transmet à la goroutine de diffusion,
 * et gère les erreurs de connexion
 */
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	var msg Message

	// Attendez un message d'inscription du client
	err = conn.ReadJSON(&msg)
	if err != nil {
		log.Printf("Error reading registration message: %v", err)
		return
	}

	username := msg.Username

	mutex.Lock()
	clients[conn] = username
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

		fmt.Printf("connections.go - %s: %s\n", msg.Username, msg.Content)

		broadcast <- msg
	}
}
