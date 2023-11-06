package server

import (
	"log"
	"fmt"
)

/**
 * Diffuse les messages aux clients via la carte des clients,
 * boucle en continu pour recevoir les messages de la goroutine handleConnections
 */
func HandleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client, username := range clients {
			// Exclure le client actuel de la diffusion
			if client != nil {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("Error writing message: %v", err)
				}
				fmt.Printf("Message envoyé à %s: %s\n", username, msg.Content)
			}
		}
		mutex.Unlock()
	}
}
