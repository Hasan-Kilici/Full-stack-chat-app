package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name            string         `gorm:"size:255;not null"`
	Email           string         `gorm:"size:255;uniqueIndex;not null"`
	Password        string         `gorm:"size:255;not null"`
	RememberToken   string         `gorm:"size:1000"`
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type PasswordReset struct {
	Email     string    `gorm:"primaryKey;size:255;index"`
	Token     string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type EmailVerification struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Token     string    `gorm:"size:255;unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Token     string    `gorm:"size:255;unique;not null"`
	UserAgent string    `gorm:"size:255"`
	IP        string    `gorm:"size:100"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type PersonalAccessToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;index;not null"`
	Token     string     `gorm:"size:64;uniqueIndex;not null"`
	Name      string     `gorm:"size:255"`
	Abilities string     `gorm:"type:text"`
	LastUsed  *time.Time
	ExpiresAt *time.Time
	CreatedAt time.Time
}

type Room struct {
	ID           uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string            `gorm:"size:255"`
	IsGroup      bool              `gorm:"default:false"`
	CreatedByID  uuid.UUID         `gorm:"type:uuid;not null"`
	CreatedBy    User              `gorm:"foreignKey:CreatedByID"`
	Participants []Participant `gorm:"foreignKey:RoomID"`
	Messages     []Message         `gorm:"foreignKey:RoomID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Participant struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RoomID uuid.UUID `gorm:"type:uuid;index;not null"`
	UserID uuid.UUID `gorm:"type:uuid;index;not null"`

	Room Room `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Message struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RoomID    uuid.UUID      `gorm:"type:uuid;index;not null"`
	UserID    uuid.UUID      `gorm:"type:uuid;index;not null"`
	Content   string         `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserID"`
	Room Room `gorm:"foreignKey:RoomID"`
}
