package server

import (
	"time"
	"fmt"
	"log"
	"net/http"
)

/**
 * Gestionnaire de connexion WebSocket
 */
func HandleConnections(w http.ResponseWriter, r *http.Request) {

    // Mise à niveau de la demande HTTP en une connexion WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
        return
    }
    // Fermeture de la connexion lorsque la fonction retourne
    defer conn.Close()

    // Enregistrement du nouveau client dans la carte
    var msg Message

    // Attente d'un message d'inscription du client
    err = conn.ReadJSON(&msg)
    if err != nil {
        log.Printf("Error reading registration message: %v", err)
        return
    }

    // Enregistrement du nouveau client dans la carte
    username := msg.Username

    fmt.Printf("%s s'est connecté.\n", username)

    // Envoie d'un message de notification "join" aux autres clients
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

    // Envoie de l'historique des messages au nouveau client
    sendChatHistory(conn)

    mutex.Lock()
    clients[conn] = username
    mutex.Unlock()

    // Boucle infinie pour lire les messages du client
    for {

        var msg Message

        err := conn.ReadJSON(&msg)

        if err != nil {
            log.Printf("Error reading message: %v", err)

            mutex.Lock()
            delete(clients, conn)
            mutex.Unlock()

            // Envoie d'un message de notification "leave" aux autres clients
            notificationMsg := Message{
                Username: username,
                Content:  "s'est déconnecté.",
                Type:     MessageTypeLeave,
                Timestamp: time.Now(),
            }
            broadcast <- notificationMsg

            break
        }

        // Ajout d'un horodatage au message reçu
        msg.Timestamp = time.Now()

        // Diffusion du message à tous les clients
        if msg.Type == MessageTypeNormal {
            fmt.Printf("connections.go [normal] - %s: %s\n", msg.Username, msg.Content)
            broadcast <- msg // Diffuser le message
            // Enregistrement du message dans l'historique
            saveToChatHistory(msg)
        }
    }
}

/**
 * Ajout du message en paramètre à l'historique des messages
 */
func saveToChatHistory(msg Message) {
    mutex.Lock()
    chatHistory = append(chatHistory, msg)
    if len(chatHistory) > 100 {
        chatHistory = chatHistory[len(chatHistory)-100:]
    }
    mutex.Unlock()
}
