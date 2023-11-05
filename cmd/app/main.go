package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/Bobybyk/go-chat-server/server"
    "github.com/gorilla/handlers"
)

/** Configuration du gestionnaire d'itinéraire pour la connexion WebSocket, 
 * démarrage d'une goroutine pour gérer les messages entrants, 
 * sert les fichiers statiques ou l'interface utilisateur (plus tard), 
 * écoute sur le port 8081
 */
func main() {
    // Créez un routeur avec gorilla/mux
    r := http.NewServeMux()

    // Configuration de CORS
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

    // Ajoutez un gestionnaire WebSocket avec le routeur
    r.HandleFunc("/ws", server.HandleConnections)

    // Créez un gestionnaire CORS avec les options configurées
    corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(r)

    // Utilisez le gestionnaire CORS comme gestionnaire racine
    http.Handle("/", corsHandler)

    // Démarrage du serveur sur le port 8080
    fmt.Println("Serveur WebSocket démarré sur le port :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
