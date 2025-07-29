package ws

import (
	"log"
	"github.com/google/uuid"
)

func writeError(c *Client, message string) {
	log.Println("Error:", message)
	_ = c.Conn.WriteJSON(Message{
		Event: "error",
		Content:  message,
	})
}

func parseUUID(idStr string, c *Client, errMsg string) (uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("%s: %v", errMsg, err)
		writeError(c, errMsg)
		return uuid.Nil, err
	}
	return id, nil
}