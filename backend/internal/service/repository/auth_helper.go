package repository

import (
	"errors"
	"time"

	"github.com/hasan-kilici/chat/internal/service/model"
	"github.com/hasan-kilici/chat/pkg/auth"
	"gorm.io/gorm"
)

func RegisterUser(name, email, password string, userAgent, ip string, rememberMe bool) (*model.User, *model.Session, error) {
    var count int64
    if err := db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
        return nil, nil, err
    }
    if count > 0 {
        return nil, nil, errors.New("email already registered")
    }

    hashedPwd, err := auth.HashPassword(password)
    if err != nil {
        return nil, nil, err
    }

    user := &model.User{
        Name:     name,
        Email:    email,
        Password: hashedPwd,
    }

    if err := db.Create(user).Error; err != nil {
        return nil, nil, err
    }

    var session *model.Session
    if rememberMe {
        token := auth.GenerateRandomToken(64)
        expires := time.Now().Add(30 * 24 * time.Hour)

        session = &model.Session{
            UserID:    user.ID,
            Token:     token,
            UserAgent: userAgent,
            IP:        ip,
            ExpiresAt: expires,
        }

        if err := db.Create(session).Error; err != nil {
            return user, nil, err
        }

        user.RememberToken = token
        if err := db.Save(user).Error; err != nil {
            return user, nil, err
        }
    }

    return user, session, nil
}

func LoginUser(email, password string) (string, *model.User, error) {
    var user model.User
    err := db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return "", nil, errors.New("user not found")
        }
        return "", nil, err
    }

    if err := auth.CheckPassword(user.Password, password); err != nil {
        return "", nil, errors.New("invalid credentials")
    }

    session := &model.Session{
        UserID:    user.ID,
        Token:     auth.GenerateRandomToken(64),
        ExpiresAt: time.Now().AddDate(0, 0, 1),
    }
    if err := db.Create(session).Error; err != nil {
        return "", nil, err
    }

    token, err := auth.GenerateJWT(user.ID, session.ID, 365*24*time.Hour)
    if err != nil {
        return "", nil, err
    }

    return token, &user, nil
}

func GetSessionByID(sessionID interface{}) (*model.Session, error) {
    var session model.Session
    if err := db.First(&session, "id = ?", sessionID).Error; err != nil {
        return nil, err
    }
    return &session, nil
}