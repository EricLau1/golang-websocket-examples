package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type Chatter struct {
	socket *websocket.Conn
	send   chan []byte
	room   *Room
}

func (c *Chatter) read() {

	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			log.Println("Erro ao ler mensagem:", err)

			break
		}
	}
	c.socket.Close()
}

func (c *Chatter) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Erro ao escrever mensagem:", err)
			break
		}
	}
	c.socket.Close()
}
