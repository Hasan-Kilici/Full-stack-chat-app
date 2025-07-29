package repository

import (
    "context"
    "encoding/json"
    "log"
    "time"

	"github.com/hasan-kilici/chat/internal/service/model"
	"github.com/google/uuid"
)

func userCacheKey(userID uuid.UUID) string {return "user:" + userID.String()}

func GetUserByID(userID uuid.UUID) (*model.User, error) {
    var user model.User
    cacheKey := userCacheKey(userID)

    val, err := redisClient.Get(ctx, cacheKey).Result()
    if err == nil && val != "" {
        if err := json.Unmarshal([]byte(val), &user); err == nil {
            return &user, nil
        }
    }

    if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
        return nil, err
    }

    if data, err := json.Marshal(user); err == nil {
        err := redisClient.Set(ctx, cacheKey, data, 10*time.Minute).Err()
        if err != nil {
            log.Printf("failed to cache user: %v", err)
        }
    }

    return &user, nil
}

func SearchUser(input string) ([]model.SafeUser, error) {
	var users []model.SafeUser

	err := db.Raw(`
		SELECT id, name, email, created_at FROM users 
		WHERE levenshtein(name, ?) < 3
		ORDER BY levenshtein(name, ?) ASC
	`, input, input).Scan(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserParticipants(ctx context.Context, userID uuid.UUID) ([]model.Participant, error) {
	var participants []model.Participant

	err := db.WithContext(ctx).
		Preload("Room").
		Preload("Room.Participants").
		Preload("Room.Participants.User").
		Where("user_id = ?", userID).
		Find(&participants).Error

	if err != nil {
		return nil, err
	}

	return participants, nil
}