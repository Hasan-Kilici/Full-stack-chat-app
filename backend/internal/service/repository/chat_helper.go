package repository

import (
	"context"
	"log"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
	"github.com/hasan-kilici/chat/internal/service/model"
)

const (
	roomCacheTTL      = 10 * time.Minute
	messagesCacheTTL  = 1 * time.Minute
	maxCachedMessages = 50
)

func roomCacheKey(roomID uuid.UUID) string    { return "room:" + roomID.String() }
func messagesCacheKey(roomID uuid.UUID) string { return "room:" + roomID.String() + ":messages" }

// CreateRoom creates a room and stores it in Redis
func CreateRoom(ctx context.Context, name string, createdBy uuid.UUID, isGroup bool, otherUserID *uuid.UUID) (*model.Room, error) {
    if !isGroup && otherUserID != nil && *otherUserID != createdBy {
        var existingRoom model.Room

        subQuery := db.Model(&model.Participant{}).
            Select("room_id").
            Group("room_id").
            Having("COUNT(*) = 2 AND BOOL_AND(user_id IN (?, ?))", createdBy, *otherUserID)

        if err := db.Where("id IN (?) AND is_group = FALSE", subQuery).
            First(&existingRoom).Error; err == nil {
            return &existingRoom, nil
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("failed to check existing dm room: %w", err)
        }
    }

    room := &model.Room{
        Name:        name,
        IsGroup:     isGroup,
        CreatedByID: createdBy,
        CreatedAt:   time.Now(),
    }

    if err := db.Create(room).Error; err != nil {
        return nil, fmt.Errorf("failed to create room: %w", err)
    }

    err := AddParticipant(ctx, room.ID, createdBy)
    if err != nil {
        return nil, fmt.Errorf("failed to add creator as participant: %w", err)
    }

    if !isGroup && otherUserID != nil && *otherUserID != createdBy {
        err := AddParticipant(ctx, room.ID, *otherUserID)
        if err != nil {
            return nil, fmt.Errorf("failed to add other user as participant: %w", err)
        }
    }

    if data, err := json.Marshal(room); err == nil {
        if err := redisClient.Set(ctx, roomCacheKey(room.ID), data, roomCacheTTL).Err(); err != nil {
            log.Printf("failed to cache room: %v", err)
        }
    }

    return room, nil
}

// GetRoom fetches a room from Redis or DB
func GetRoom(ctx context.Context, roomID uuid.UUID) (*model.Room, error) {
	var room model.Room

	if val, err := redisClient.Get(ctx, roomCacheKey(roomID)).Result(); err == nil && val != "" {
		if jsonErr := json.Unmarshal([]byte(val), &room); jsonErr == nil {
			return &room, nil
		}
	}

	if err := db.Preload("Participants").First(&room, "id = ?", roomID).Error; err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}

	if data, err := json.Marshal(room); err == nil {
		err := redisClient.Set(ctx, roomCacheKey(roomID), data, roomCacheTTL).Err()
		if err != nil {
			log.Printf("failed to cache room: %v", err)
		}
	}

	return &room, nil
}

func GetMessages(ctx context.Context, roomID uuid.UUID, limit, page int) ([]model.Message, error) {
	var messages []model.Message

	offset := (page - 1) * limit
	cacheKey := fmt.Sprintf("%s:limit=%d:page=%d", messagesCacheKey(roomID), limit, page)

	if val, err := redisClient.Get(ctx, cacheKey).Result(); err == nil && val != "" {
		if jsonErr := json.Unmarshal([]byte(val), &messages); jsonErr == nil {
			return messages, nil
		}
	}

	if err := db.
		Where("room_id = ?", roomID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	if data, err := json.Marshal(messages); err == nil {
		_ = redisClient.Set(ctx, cacheKey, data, messagesCacheTTL).Err()
	}

	return messages, nil
}

func GetSafeMessages(ctx context.Context, roomID uuid.UUID, limit, page int) ([]model.SafeMessage, error) {
	offset := (page - 1) * limit
	cacheKey := fmt.Sprintf("%s:limit=%d:page=%d", messagesCacheKey(roomID), limit, page)

	var safeMessages []model.SafeMessage
	if val, err := redisClient.Get(ctx, cacheKey).Result(); err == nil && val != "" {
		if jsonErr := json.Unmarshal([]byte(val), &safeMessages); jsonErr == nil {
			return safeMessages, nil
		}
	}

	var messages []model.Message
	if err := db.
		Preload("User").
		Where("room_id = ?", roomID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	safeMessages = make([]model.SafeMessage, 0, len(messages))
	for _, msg := range messages {
		safeMessages = append(safeMessages, model.SafeMessage{
			ID:        msg.ID,
			AuthorID:  msg.UserID,
			Author:    msg.User.Name,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		})
	}

	if data, err := json.Marshal(safeMessages); err == nil {
		_ = redisClient.Set(ctx, cacheKey, data, messagesCacheTTL).Err()
	}

	return safeMessages, nil
}


func SaveMessage(ctx context.Context, roomID, userID uuid.UUID, content string) error {
	msg := &model.Message{
		RoomID:    roomID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := db.Create(msg).Error; err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	messages, err := GetMessages(ctx, roomID, maxCachedMessages, 0)
	if err != nil {
		messages = []model.Message{}
	}

	messages = append([]model.Message{*msg}, messages...)
	if len(messages) > maxCachedMessages {
		messages = messages[:maxCachedMessages]
	}

	return nil
}

func AddParticipant(ctx context.Context, roomID, userID uuid.UUID) error {
    var count int64
    db.WithContext(ctx).Model(&model.Participant{}).
        Where("room_id = ? AND user_id = ?", roomID, userID).
        Count(&count)
    if count > 0 {
        return nil
    }
    participant := &model.Participant{
        RoomID: roomID,
        UserID: userID,
    }
    if err := db.WithContext(ctx).Create(participant).Error; err != nil {
        return fmt.Errorf("failed to add participant: %w", err)
    }
    return nil
}

func GetParticipants(ctx context.Context, roomID uuid.UUID) ([]model.Participant, error) {
	var participants []model.Participant
	if err := db.WithContext(ctx).Where("room_id = ?", roomID).Find(&participants).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch participants: %w", err)
	}
	return participants, nil
}
