package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	//time allowed to write a message to the peer
	writeWait = 10 * time.Second

	//time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	//send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	//maximum message size allowed from peer
	maxMessageSize = 512
)

// Client represents a single WebSocket connection
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte //buffered channel
	room     string
	username string
}

// pumps messages from the WebSocket to the Hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		//parse incoming message
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("error parsing message: %v", err)
			continue
		}

		msg.Username = c.username

		//handle different message types
		switch msg.Type {
		case "join":
			skipBroadcast := false
			if msg.SkipBroadcast {
				skipBroadcast = msg.SkipBroadcast
			}
			c.hub.joinRoom <- &RoomRequest{client: c, room: msg.Room, skipBroadcast: skipBroadcast}
		case "leave":
			c.hub.leaveRoom <- &RoomRequest{client: c, room: msg.Room}
		case "message":
			if c.room != "" {
				msg.Room = c.room
				c.hub.broadcast <- &msg
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			//drain queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
