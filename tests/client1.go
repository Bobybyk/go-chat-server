package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Message struct {
	Content string `json:"content"`
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	serverURL := "ws://localhost:8080/ws"

	c, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Erreur de connexion au serveur WebSocket:", err)
	}
	defer c.Close()

	done := make(chan struct{})

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
			fmt.Printf("Message reçu: %s\n", msg.Content)
		}
	}()

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
			fmt.Print("Entrez un message: ")
			inputReader := bufio.NewReader(os.Stdin)
			message, _ := inputReader.ReadString('\n')
			message = message[:len(message)-1] // Suppression du caractère de nouvelle ligne
			msg := Message{Content: message}
			msgBytes, _ := json.Marshal(msg)
			err = c.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				log.Println("Erreur d'envoi du message:", err)
				return
			}
		}
	}
}
