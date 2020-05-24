package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Erro ao ler mensagem:", err)
			return
		}

		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Erro ao escrever mensagem:", err)
			return
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Cliente Conectado com Sucesso!")

	reader(ws)

}

func setupRoutes() {
	http.HandleFunc("/", checkHealth)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Go Websocket Example")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
