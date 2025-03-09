package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub		*Hub
	Conn	*websocket.Conn
	UserID	int64
	MatchID	int64

	// channel for sending message.
	send chan BroadcastMessage
}

// ReadPump continuously reads messages from the client's WebSocket connection.
func (c *Client) ReadPump(handleMessage func(c *Client, msg string) error) {
	// Ensure the client is unregistered and the connection is closed when the function exits.
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	// Loop indefinitely to continuously read messages.
	for {
		// Read a message from the WebSocket connection.
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v\n", err)
			break
		}
		// Process the message using the provided handleMessage function.
		if err := handleMessage(c, string(message)); err != nil {
			log.Printf("handleMessage error: %v\n", err)
		}
	}
}

func (c *Client) WritePump() {
	// Ensure the WebSocket connection is closed when this function exits.
	defer func() {
		c.Conn.Close()
	}()

	// Loop indefinitely to process messages from the send channel.
	for {
		select {
		// Wait for a message from the send channel.
		case message, ok := <-c.send:
			// If the channel is closed, send a WebSocket close message and exit the loop.
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// Write the message as JSON to the WebSocket connection.
			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("write error: %v\n", err)
				return
			}
		}
	}
}