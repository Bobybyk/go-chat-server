package server

import (
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
