package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/hasan-kilici/chat/internal/service/repository"
)

type CreateRoomRequest struct {
	Name        string     `json:"name"`                      
	IsGroup     bool       `json:"is_group"`                 
	OtherUserID *uuid.UUID `json:"other_user_id,omitempty"`  
}

func CreateRoomHandler(c fiber.Ctx) error {
    userIDStr := c.Locals("user_id")
    if userIDStr == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Unauthorized",
        })
    }
    userID, err := uuid.Parse(userIDStr.(string))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid user ID",
        })
    }

    var req CreateRoomRequest
    if err := json.Unmarshal(c.Body(), &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid JSON body",
        })
    }

    if !req.IsGroup {
        if req.OtherUserID == nil || *req.OtherUserID == userID {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "other_user_id is required and must be different from your ID",
            })
        }
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    room, err := repository.CreateRoom(ctx, req.Name, userID, req.IsGroup, req.OtherUserID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not create room",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "room_id":    room.ID,
        "is_group":   room.IsGroup,
        "created_by": room.CreatedByID,
        "created_at": room.CreatedAt,
    })
}

func JoinRoomHandler(c fiber.Ctx) error {
    userIDStr := c.Locals("user_id")
    if userIDStr == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }
    userID, err := uuid.Parse(userIDStr.(string))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
    }

    roomIDStr := c.Params("room_id")
    roomID, err := uuid.Parse(roomIDStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid room_id"})
    }

    room, err := repository.GetRoom(context.Background(), roomID)
    if err != nil || room == nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Room not found"})
    }

    if !room.IsGroup && len(room.Participants) >= 2 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot join a private room with more than one participant"})
    }

    if err := repository.AddParticipant(context.Background(), roomID, userID); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Joined room successfully"})
}
