package main

import (
    "fmt"
    "log"
    // Importation du gestionnaire de route
    "net/http"
    // Importation du gestionnaire de route
    "github.com/Bobybyk/go-chat-server/server"
    // Importation du gestionnaire CORS
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
 * Possible d'ajuster des options CORS en fonction des besoins en sécurité et de gestion des origines.
 */
func main() {
    // Créez un routeur avec gorilla/mux pour gérer les routes
    r := http.NewServeMux()

    // Configuration de CORS

    // Autorise les en-têtes X-Requested-With et Content-Type
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
    // Autorise les requêtes depuis n'importe quelle origine
    originsOk := handlers.AllowedOrigins([]string{"*"})
    // Autorise les méthodes GET, POST, PUT, OPTIONS
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

    // Ajout d'un gestionnaire WebSocket avec le routeur
    r.HandleFunc("/ws", server.HandleConnections)

    // Cration d'un gestionnaire CORS avec les options configurées
    corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(r)

    // Gestionnaire CORS utilisé comme gestionnaire racine
    http.Handle("/", corsHandler)

    // Démarrage d'une goroutine pour gérer les messages entrants
    go server.HandleMessages()
    
    // Démarrage du serveur sur le port 8080
    fmt.Println("Serveur WebSocket démarré sur le port :8080")
    err := http.ListenAndServe(":55396", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
