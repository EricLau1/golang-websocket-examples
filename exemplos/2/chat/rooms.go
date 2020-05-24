package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type Room struct {
	forward  chan []byte
	join     chan *Chatter
	leave    chan *Chatter
	chatters map[*Chatter]bool
	topic    string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("Check origin: %s %s%s", r.Method, r.Host, r.URL.Path)
		return true
	},
}

func (r *Room) Handler(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("Handler HTTP falhou:", err)
	}

	chatter := &Chatter{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- chatter

	defer func() {
		r.leave <- chatter
	}()

	go chatter.write()
	chatter.read()
}

func NewRoom(topic string) *Room {
	return &Room{
		forward:  make(chan []byte),
		join:     make(chan *Chatter),
		leave:    make(chan *Chatter),
		chatters: make(map[*Chatter]bool),
		topic:    topic,
	}
}

func (r *Room) Run() {
	log.Printf("Chat/%s\n", r.topic)
	for {
		select {
		case chatter := <-r.join:
			log.Printf("Alguém se juntou a sala: %v", r.topic)
			r.chatters[chatter] = true

		case chatter := <-r.leave:
			log.Printf("Alguém saiu da sala: %v", r.topic)
			delete(r.chatters, chatter)
			close(chatter.send)

		case input := <-r.forward:
			msg := fromJSON(input)
			log.Printf("%s/%s escreveu: %v", r.topic, msg.Sender, msg.Body)
			r.Broadcast(input)
		}
	}
}

func (r *Room) Broadcast(msg []byte) {
	for chatter := range r.chatters {
		select {
		case chatter.send <- msg:
		default:
			delete(r.chatters, chatter)
			close(chatter.send)
		}
	}
}
