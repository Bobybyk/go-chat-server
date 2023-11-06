package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/Bobybyk/go-chat-server/server"
    "github.com/gorilla/handlers"
)

/** Configuration du gestionnaire d'itinéraire pour la connexion WebSocket, 
 * configuration du gestionnaire CORS,
 * démarrage d'une goroutine pour gérer les messages entrants, 
 * sert les fichiers statiques ou l'interface utilisateur (plus tard), 
 * écoute sur le port 8080
 * 
 * Note sur le gestionnaire CORS:
 * Toutes les demandes HTTP et WebSocket reçues par le serveur passeront d'abord par le gestionnaire CORS, 
 * qui vérifiera si elles sont autorisées en fonction des options spécifiées. 
 * Si une demande correspond aux critères définis (origine, méthodes, en-têtes), elle sera autorisée.
 * 
 * Avec cette configuration, le serveur acceptera des connexions WebSocket depuis n'importe quel domaine, 
 * ce qui facilite la connexion depuis d'autres machines ou sites web. 
 * Possible d'ajuster les options CORS en fonction des besoins en sécurité et de gestion des origines.
 */
func main() {
    // Créez un routeur avec gorilla/mux pour gérer les routes
    r := http.NewServeMux()

    // Configuration de CORS
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

    // Ajout d'un gestionnaire WebSocket avec le routeur
    r.HandleFunc("/ws", server.HandleConnections)

    // Cration d'un gestionnaire CORS avec les options configurées
    corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(r)

    // Gestionnaire CORS utilisé comme gestionnaire racine
    http.Handle("/", corsHandler)
    go server.HandleMessages()
    
    // Démarrage du serveur sur le port 8080
    fmt.Println("Serveur WebSocket démarré sur le port :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
