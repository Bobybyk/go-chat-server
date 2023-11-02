package server

import (
	"log"
)

/**
 * Diffuse les messages aux clients via la carte des clients,
 * boucle en continu pour recevoir les messages de la goroutine handleConnections
 */
func HandleMessages() {
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
