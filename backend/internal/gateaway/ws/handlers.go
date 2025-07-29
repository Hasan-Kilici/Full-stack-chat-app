package ws

import (
	"fmt"
	"log"
	"net/http"
	"context"

	"github.com/hasan-kilici/chat/internal/service/repository"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)


func init() {
	Router.Register("join", handleJoin)
	Router.Register("message", handleMessage)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	authResult, ok := checkJWTAuth(r, w)
	if !ok {
		return
	}

	fmt.Println("WebSocket connection established for user:", authResult)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrader error", err)
		return
	}

	var client *Client

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			fmt.Println("connection over", err)
			if client != nil {
				Hub.Leave(client)
			}
			break
		}

		if client == nil && msg.Event == "join" {
			client = &Client{
				Conn: conn,
				User: User{
					ID:   authResult.UserID,
					Name: authResult.Name,
				},
				Room: msg.Room,
			}
			fmt.Println("New client connected:", client.User.Name, "to room:", client.Room)
			Hub.Join(client)
		}

		Router.Handle(client, msg)
	}
}

func handleJoin(c *Client, msg Message) {
	var canJoin = false
	roomID, err := parseUUID(msg.Room, c, "Invalid room ID.")
	if err != nil {
		return
	}

	room, err := repository.GetRoom(context.Background(), roomID)
	if err != nil || room == nil {
		writeError(c, "Room not found.")
		return
	}

	userID, err := uuid.Parse(c.User.ID)
	if err != nil {
		writeError(c, "Invalid user ID.")
		return
	}

	for _, participant := range room.Participants {
		if participant.UserID == userID {
			fmt.Printf("📝 %s sent a message in room %s: %s\n", c.User.Name, room.ID, msg.Content)
			fmt.Println("User already in room:", c.User.Name, "Room ID:", room.ID)
			canJoin = true
			break
		}
	}

	if !room.IsGroup && len(room.Participants) >= 2 && !canJoin {
		writeError(c, "Cannot join a private room with more than one participant.")
		return
	}

	if err := repository.AddParticipant(context.Background(), room.ID, userID); err != nil {
		writeError(c, "Could not join the room.")
		return
	}

	fmt.Printf("🔗 %s joined the room: %s\n", c.User.Name, room.ID)

	Hub.Broadcast(c, Message{
		Event: 		"join",
		Author:  	c.User.Name,
		AuthorID: 	c.User.ID,
		Room:  		room.Name,
		Content:  	"Joined the room.",
	})
}

func handleMessage(c *Client, msg Message) {
	var canMessage = false
	if msg.Room == "" {
		writeError(c, "Room is required.")
		return
	}
	if msg.Content == "" {
		writeError(c, "Message data cannot be empty.")
		return
	}

	roomID, err := parseUUID(msg.Room, c, "Invalid room ID format.")
	if err != nil {
		return
	}

	room, err := repository.GetRoom(context.Background(), roomID)
	if err != nil || room == nil {
		writeError(c, "Room not found.")
		return
	}

	userID, err := uuid.Parse(c.User.ID)
	if err != nil {
		writeError(c, "Invalid user ID.")
		return
	}

	for _, participant := range room.Participants {
		if participant.UserID == userID {
			fmt.Printf("📝 %s sent a message in room %s: %s\n", c.User.Name, room.ID, msg.Content)
			fmt.Println("User already in room:", c.User.Name, "Room ID:", room.ID)
			canMessage = true
			break
		}
	}

	if !room.IsGroup && len(room.Participants) >= 2 && !canMessage{
		writeError(c, "Cannot send messages in a private room with more than one participant.")
		return
	}

	if err := repository.SaveMessage(context.Background(), roomID, userID, msg.Content); err != nil {
		writeError(c, "Could not send message.")
		return
	}

	Hub.Broadcast(c, Message{
		Event: 		msg.Event,
		Author:  	c.User.Name,
		AuthorID:	c.User.ID,	
		Room:  		msg.Room,
		Content:  	msg.Content,
	})
}