package main

import (
    "log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,

    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }

    defer ws.Close()

    log.Println("New Client Connected")

    for {
        messageType, p, err := ws.ReadMessage()
        if err != nil {
            log.Println(err)
			return
        }

        log.Printf("Message Received: %s\n", p)

        if err := ws.WriteMessage(messageType, p); err != nil {
            log.Println(err)
			return
        }
    }
}

func main() {
    http.HandleFunc("/ws", handleConnections)

    log.Println("http server started on :8080")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}