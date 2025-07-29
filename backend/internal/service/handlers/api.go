package handlers

import (
	"strconv"
	"github.com/hasan-kilici/chat/internal/service/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func SearchUser(c fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	users, err := repository.SearchUser(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search users",
		})
	}

	return c.JSON(users)
}

func GetMessages(c fiber.Ctx) error {
	roomID := c.Params("room_id")
	if roomID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Room ID is required",
		})
	}

	page := c.Query("page", "1")
	pageSize := c.Query("page_size", "100")
	pageInt, err := strconv.Atoi(page)

	if err != nil || pageInt < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page number",
		})
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page size",
		})
	}

	roomIDParsed, err := uuid.Parse(roomID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Room ID format",
		})
	}

	messages, err := repository.GetSafeMessages(c, roomIDParsed, 100, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve messages",
		})
	}

	return c.JSON(messages)
}

func GetUserRooms(c fiber.Ctx) error {
	userIDstr := c.Locals("user_id")
	userID, err := uuid.Parse(userIDstr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	participants, err := repository.GetUserParticipants(c, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve rooms"})
	}

	type RoomResponse struct {
		RoomID   uuid.UUID `json:"room_id"`
		RoomName string    `json:"room_name"`
	}

	var result []RoomResponse

	for _, p := range participants {
		room := p.Room
		roomName := ""

		if room.IsGroup {
			roomName = room.Name
		} else {
			for _, otherP := range room.Participants {
				if otherP.UserID != userID {
					roomName = otherP.User.Name
					break
				}
			}
		}

		result = append(result, RoomResponse{
			RoomID:   room.ID,
			RoomName: roomName,
		})
	}

	return c.JSON(result)
}