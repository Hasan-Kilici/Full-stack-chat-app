package model

import (
	"time"
	"github.com/google/uuid"
)

type (
	SafeUser struct {
		ID        	uuid.UUID 	`json:"id"`
		Name      	string    	`json:"name"`
		Email     	string    	`json:"email"`
		CreatedAt 	time.Time 	`json:"created_at"`
	}

	SafeParticipant struct {
		RoomID    	uuid.UUID 	`json:"room_id"`
		RoomName  	string    	`json:"room_name"`
	}

	SafeMessage struct {
		ID        uuid.UUID `json:"id"`
		AuthorID  uuid.UUID `json:"author_id"`
		Author    string    `json:"author"`
		Content   string    `json:"data"`
		CreatedAt time.Time `json:"created_at"`
	}
)