package chat

import (
	"encoding/json"
	"log"
)

type Message struct {
	Body    string `json:"body"`
	Sender  string `json:"sender"`
	Created string `json:"created"`
}

func fromJSON(input []byte) (message *Message) {
	if err := json.Unmarshal(input, &message); err != nil {
		log.Println("Erro ao deserializar json input:", err)
	}
	return
}
