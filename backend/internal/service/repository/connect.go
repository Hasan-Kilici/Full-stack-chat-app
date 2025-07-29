package repository

import (
	"context"
	"log"

	"github.com/hasan-kilici/chat/pkg/utils"
	"github.com/hasan-kilici/chat/internal/service/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	redisClient *redis.Client
	ctx       = context.Background()
)

func Connect() {
	config := utils.LoadConfig("./configs/service.ini")

	var err error
	db, err = gorm.Open(postgres.Open(config.Database), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.PasswordReset{},
		&model.EmailVerification{},
		&model.Session{},
		&model.PersonalAccessToken{},
		&model.Room{},
		&model.Participant{},
		&model.Message{},
	)
	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}
	log.Println("Database migration completed")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		DB:       0,
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}