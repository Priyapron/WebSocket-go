package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		conn.Close()
		delete(clients, conn)
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		message := string(p)

		// ทำการ broadcast ข้อความที่ได้รับจาก client ไปทุก client ที่เชื่อมต่อ
		clients[conn] = true
		for client := range clients {
			if client != conn {
				err := client.WriteMessage(messageType, []byte(message))
				if err != nil {
					log.Printf("Error writing to client: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

var clients = make(map[*websocket.Conn]bool)

func main() {
	http.HandleFunc("/ws", handleConnections)

	// บริการไฟล์ HTML แยกออกมา
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	fmt.Println("Server is running on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
