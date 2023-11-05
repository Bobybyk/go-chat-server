package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"os/signal"
	"encoding/json"
	"github.com/gorilla/websocket"
)

// Structure du message échangé
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Type     string `json:"type"`
}

const (
	MessageTypeNormal    = "normal"
	MessageTypeJoin      = "join"
	MessageTypeLeave     = "leave"
)


func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	serverURL := "ws://localhost:8081/ws" // Mettez l'URL de votre serveur ici

	c, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Erreur de connexion au serveur WebSocket:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Saisie du nom d'utilisateur
	var username string
	fmt.Print("Entrez votre nom d'utilisateur: ")
	_, err = fmt.Scanln(&username)
	if err != nil {
		log.Println("Erreur de lecture de l'entrée standard:", err)
		return
	}

	// Envoyer un message d'inscription au serveur
	registerMessage := Message{Username: username, Content: ""}
	registerMessageBytes, _ := json.Marshal(registerMessage)
	err = c.WriteMessage(websocket.TextMessage, registerMessageBytes)
	if err != nil {
		log.Println("Erreur d'envoi du message d'inscription:", err)
		return
	}

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Erreur de lecture du message:", err)
				return
			}
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("Erreur de désérialisation du message:", err)
				return
			}
			if msg.Type == MessageTypeJoin {
				fmt.Printf("%s %s\n", msg.Username, msg.Content)
			} else if msg.Type == MessageTypeLeave {
				fmt.Printf("%s %s\n", msg.Username, msg.Content)
			} else {
				fmt.Printf("%s: %s\n", msg.Username, msg.Content)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Vous êtes connecté (Ctrl+Z pour quitter):")
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("Fermeture de la connexion...")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Erreur de fermeture de la connexion:", err)
				return
			}
			select {
			case <-done:
			}
			return
		default:
			scanner.Scan()
			message := scanner.Text()
			msg := Message{Username: username, Content: message}
			msgBytes, _ := json.Marshal(msg)
			err = c.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				log.Println("Erreur d'envoi du message:", err)
				return
			}
		}
	}
}
