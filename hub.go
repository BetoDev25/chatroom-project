package main

import (
	"encoding/json"
	"log"
)

type RoomRequest struct {
	client        *Client
	room          string
	skipBroadcast bool
}

type JoinMessage struct {
	Type          string `json:"type"`
	Room          string `json:"room"`
	SkipBroadcast bool   `json:"skipBroadcast,omitempty"`
}

// Hub maintains active clients and broadcasts messages to rooms
type Hub struct {
	//all connected clients
	clients map[*Client]bool

	//room name to clients mapping
	rooms map[string]map[*Client]bool

	//convo name to clients mapping
	conversations map[string]map[*Client]bool

	//inbound messages to broadcast
	broadcast chan *Message

	//register requests from clients
	register chan *Client

	//unregister requests from clients
	unregister chan *Client

	//join room requests
	joinRoom chan *RoomRequest

	//leave room requests
	leaveRoom chan *RoomRequest
}

func newHub() *Hub {
	return &Hub{
		broadcast:     make(chan *Message),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		joinRoom:      make(chan *RoomRequest),
		leaveRoom:     make(chan *RoomRequest),
		clients:       make(map[*Client]bool),
		rooms:         make(map[string]map[*Client]bool),
		conversations: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected: %s", client.username)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				//remove from current room if in one
				if client.room != "" {
					h.removeFromRoom(client, client.room)
				}
				if client.conversationID != "" {
					h.leaveConversation(client)
				}
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected: %s", client.username)
			}

		case req := <-h.joinRoom:
			h.handleJoinRoom(req.client, req.room, req.skipBroadcast)

		case req := <-h.leaveRoom:
			h.removeFromRoom(req.client, req.room)

		case msg := <-h.broadcast:
			h.broadcastToRoom(msg)
		}
	}
}

func (h *Hub) handleJoinRoom(client *Client, room string, skipBroadcast bool) {
	//leave current room if in one
	if client.room != "" {
		h.removeFromRoom(client, client.room)
	}

	//create room if it doesn't exist
	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*Client]bool)
		log.Printf("Room created: %s", room)
	}

	//add client to room
	h.rooms[room][client] = true
	client.room = room
	log.Printf("%s joined room: %s", client.username, room)

	//notify room members
	if !skipBroadcast {
		joinMsg := &Message{
			Type:     "join",
			Room:     room,
			Username: client.username,
			Content:  client.username + " joined the room",
		}
		h.broadcastToRoom(joinMsg)
	}
}

func (h *Hub) removeFromRoom(client *Client, room string) {
	if clients, ok := h.rooms[room]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)

			//notify room about departure
			leaveMsg := &Message{
				Type:     "leave",
				Room:     room,
				Username: client.username,
				Content:  client.username + " left the room",
			}
			h.broadcastToRoom(leaveMsg)

			//clean up empty rooms
			if len(clients) == 0 {
				delete(h.rooms, room)
				log.Printf("Room deleted: %s", room)
			}

			client.room = ""
			log.Printf("%s left room: %s", client.username, room)
		}
	}
}

func (h *Hub) broadcastToRoom(msg *Message) {
	clients, ok := h.rooms[msg.Room]
	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshaling message: %v", err)
		return
	}

	for client := range clients {
		select {
		case client.send <- data:
			//Message sent successfully
		default:
			//Client buffer full, disconnect them
			close(client.send)
			delete(clients, client)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) joinConversation(client *Client, conversationID string) {
	if client.conversationID != "" {
		h.leaveConversation(client)
	}

	if h.conversations[conversationID] == nil {
		h.conversations[conversationID] = make(map[*Client]bool)
	}
	h.conversations[conversationID][client] = true
	client.conversationID = conversationID
}

func (h *Hub) leaveConversation(client *Client) {
	if client.conversationID == "" {
		return
	}
	if clients, ok := h.conversations[client.conversationID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.conversations, client.conversationID)
		}
		client.conversationID = ""
	}
}

func (h *Hub) broadcastPrivate(msg *Message) {
	clients, ok := h.conversations[msg.ConversationID]
	if !ok {
		log.Printf("No clients in conversation: %s", msg.ConversationID) // Debug
		return
	}

	log.Printf("Broadcasting to %d clients in conversation %s", len(clients), msg.ConversationID) // Debug

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshaling private message: %v", err)
		return
	}

	for client := range clients {
		select {
		case client.send <- data:
			log.Printf("Sent to client: %s", client.username) // Debug
		default:
			log.Printf("Client %s buffer full, disconnecting", client.username) // Debug
			close(client.send)
			delete(clients, client)
		}
	}
}
