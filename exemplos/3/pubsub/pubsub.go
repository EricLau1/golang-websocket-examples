package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

const (
	PUBLISH     = "publish"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
)

type PubSub interface {
	AddClient(c *Client)
	Dispatch(client *Client, messageType int, payload []byte)
	RemoveClient(c *Client)
}

type pubSubImpl struct {
	clients       map[string]*Client
	subscriptions map[string][]*Subscription
	lock          *sync.Mutex
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type Payload struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

type Subscription struct {
	Topic  string
	Client *Client
}

func NewPubSub() PubSub {
	return &pubSubImpl{
		clients:       make(map[string]*Client),
		subscriptions: make(map[string][]*Subscription),
		lock:          &sync.Mutex{},
	}
}

func (ps *pubSubImpl) AddClient(c *Client) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	if _, ok := ps.clients[c.ID]; !ok {
		ps.clients[c.ID] = c

		payload := []byte("Hello Client ID: " + c.ID)
		if err := c.Conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			fmt.Println("Error on Write Message: ", err.Error())
		}
	}
}

func (ps *pubSubImpl) Dispatch(client *Client, messageType int, p []byte) {

	payload := new(Payload)
	err := json.Unmarshal(p, payload)
	if err != nil {
		fmt.Println("Error on Unmarshal: ", err.Error())

		if err := client.Conn.WriteMessage(messageType, []byte("Error: "+err.Error())); err != nil {
			fmt.Println("Error on Write Message: ", err.Error())
		}
		return
	}

	ps.dispatch(client, payload, messageType)

	if err := client.Conn.WriteMessage(messageType, p); err != nil {
		fmt.Println("Error on Write Message: ", err.Error())
	}
}

func (ps *pubSubImpl) dispatch(c *Client, payload *Payload, messageType int) {
	switch payload.Action {
	case PUBLISH:
		ps.publish(payload.Topic, c, payload.Message, messageType)
	case SUBSCRIBE:
		ps.subscribe(payload.Topic, c)
		message := json.RawMessage(fmt.Sprintf("%s entered the Topic: %s", c.ID, payload.Topic))
		ps.publish(payload.Topic, c, message, messageType)
	case UNSUBSCRIBE:
		ps.unsubscribe(c.ID)
		message := json.RawMessage(fmt.Sprintf("%s left the Topic: %s", c.ID, payload.Topic))
		ps.publish(payload.Topic, c, message, messageType)
	default:
		fmt.Println("Invalid action: ", payload.Action)
	}
}

func (ps *pubSubImpl) subscribe(topic string, c *Client) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.subscriptions[topic] = append(ps.subscriptions[topic], &Subscription{Topic: topic, Client: c})
	fmt.Printf("Subscribe: %s, Topic: %s\n", c.ID, topic)
}

func (ps *pubSubImpl) publish(topic string, c *Client, message []byte, messageType int) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	fmt.Printf("Publishing Topic: %s, Message: %s\n", topic, message)
	if subscriptions, ok := ps.subscriptions[topic]; ok {
		for _, sub := range subscriptions {
			if sub.Client.ID != c.ID {
				err := sub.Client.Conn.WriteMessage(messageType, message)
				if err != nil {
					fmt.Println("Error on publish message: ", err.Error())
					continue
				}
			}
		}
	} else {
		fmt.Println("Subscriptions is empty...")
	}
}

func (ps *pubSubImpl) RemoveClient(c *Client) {
	ps.unsubscribe(c.ID)
	delete(ps.clients, c.ID)
	_ = c.Conn.Close()
}

func (ps *pubSubImpl) unsubscribe(ID string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	fmt.Printf("Unsubscribe: %s\n", ID)
	for topic, subs := range ps.subscriptions {
		for index, sub := range subs {
			if sub.Client.ID == ID {
				ps.subscriptions[topic] = append(ps.subscriptions[topic][:index], ps.subscriptions[topic][index+1:]...)
			}
		}
	}
}
