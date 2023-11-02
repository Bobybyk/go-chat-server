package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/Bobybyk/go-chat-server/server"
)

/** Configuration du gestionnaire d'itinéraire pour la connexion WebSocket, 
 * démarrage d'une goroutine pour gérer les messages entrants, 
 * sert les fichiers statiques ou l'interface utilisateur (plus tard), 
 * écoute sur le port 8081
 */
func main() {
	http.HandleFunc("/ws", server.HandleConnections)
	go server.HandleMessages()

	fmt.Println("Server started on :8082")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
