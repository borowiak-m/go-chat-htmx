package main

import (
	"bytes"
	"html/template"
	"log"

	"github.com/google/uuid"
)

type Hub struct {
	//sync.mutex
	clients              map[*Client]bool
	broadcast            chan *Message
	register, unregister chan *Client
	messages             []*Message
}

type Message struct {
	ClientID uuid.UUID
	Text     string
}

type WSMessage struct {
	Text    string      `json:"text"`
	Headers interface{} `json:"HEADERSE"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			//mutex lock
			h.clients[client] = true
			log.Printf("client registered: %s", client.id)
			//mutex unlock
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Printf("client unregistered: %s", client.id)
				//mutex lock
				close(client.send)
				delete(h.clients, client)
				//mutex unlock
			}
		case msg := <-h.broadcast:
			//mutex lock
			h.messages = append(h.messages, msg)
			for client := range h.clients {
				select {
				case client.send <- getMessageTemplate(msg):
					log.Printf("message sent: %s", msg.Text)
				default:
					//mutex lock
					close(client.send)
					delete(h.clients, client)
					//mutex unlock
				}

			}
			//mutex unlock
		}
	}
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("templates/message.html")
	if err != nil {
		log.Fatalf("error parsing files: %s", err)
	}
	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, msg)
	if err != nil {
		log.Fatalf("error rendering message: %s", err)
	}
	return rendered.Bytes()
}
