package server

import (
	"time"
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

	fmt.Printf("%s s'est connecté.\n", username)

	// Envoyer un message de notification "join" aux autres clients

	notificationMsg := Message{
		Username: username,
		Content:  "s'est connecté.",
		Type:     MessageTypeJoin,
	}
	broadcast <- notificationMsg

	// Envoie au nouveau client la liste des clients connectés
	mutex.Lock()
	for _, username := range clients {
		msg := Message{
			Username: username,
			Content:  "",
			Type:     MessageTypeJoin,
			Timestamp: time.Now(),
		}
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error writing message: %v", err)
		}
	}
	mutex.Unlock()


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
	
			// Envoyer un message de notification "leave" aux autres clients
			notificationMsg := Message{
				Username: username,
				Content:  "s'est déconnecté.",
				Type:     MessageTypeLeave,
				Timestamp: time.Now(),
			}
			broadcast <- notificationMsg
	
			break
		}
		msg.Timestamp = time.Now()
	
		if msg.Type == MessageTypeNormal {
			fmt.Printf("connections.go [normal] - %s: %s\n", msg.Username, msg.Content)
			broadcast <- msg // Diffusez le message 
		}
	}
}
