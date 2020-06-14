package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go-web-socket-examples/exemplos/3/pubsub"
	"log"
	"net/http"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin:       func(r *http.Request) bool { return true },
}

func GetID() string {
	return uuid.New().String()
}

var ps = pubsub.NewPubSub()

func wsHandler(w http.ResponseWriter, r *http.Request) {

	time.Sleep(time.Second * 3)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error on Upgrade request: ", err.Error())
		return
	}

	client := &pubsub.Client{ID: GetID(), Conn: conn }
	ps.AddClient(client)
	fmt.Println("Client connected: ", client.ID)

	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error on Read Message: ", err.Error())

			if strings.Contains(err.Error(), "close") {
				fmt.Println("Client closed session: ", client.ID)
				ps.RemoveClient(client)
			}

			return
		}

		ps.Dispatch(client, msgType, message)
	}

	if conn != nil {
		conn.Close()
	}
}

const PORT = ":8080"

func main() {
	fmt.Println("Server listening on: http://localhost:8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r * http.Request) {
		http.ServeFile(w, r, "exemplos/3/static")
	})
	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
