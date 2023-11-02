package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/Bobybyk/go-chat-server/server"
)

func main() {
	http.HandleFunc("/ws", server.HandleConnections)
	go server.HandleMessages()

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server started on :8082")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
