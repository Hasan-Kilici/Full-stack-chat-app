package handlers

import (
    "encoding/json"
    
    "github.com/gofiber/fiber/v3"
    "github.com/google/uuid"
    "github.com/hasan-kilici/chat/internal/service/repository"
)

type RegisterRequest struct {
    Name       string `json:"name"`
    Email      string `json:"email"`
    Password   string `json:"password"`
    RememberMe bool   `json:"rememberMe"`
}

func RegisterHandler(c fiber.Ctx) error {
    var req RegisterRequest
    if err := json.Unmarshal(c.Body(), &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid JSON body",
        })
    }

    user, _, err := repository.RegisterUser(
        req.Name,
        req.Email,
        req.Password,
        c.Get("User-Agent"),
        c.IP(),
        req.RememberMe,
    )

    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    token,  _, err := repository.LoginUser(req.Email, req.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "User registered successfully",
        "user_id": user.ID,
        "auth_token": token,
    })
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func LoginHandler(c fiber.Ctx) error {
    var req LoginRequest
    if err := json.Unmarshal(c.Body(), &req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid JSON body",
        })
    }

    token, user, err := repository.LoginUser(req.Email, req.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "token":   token,
        "user_id": user.ID,
        "email":   user.Email,
        "name":    user.Name,
    })
}

func AuthRequired(c fiber.Ctx) error {
    userID := c.Locals("user_id")
    if userID == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Unauthorized",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Authenticated",
        "user_id": userID,
    })
}

func GetUserProfile(c fiber.Ctx) error {
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

    user, err := repository.GetUserByID(userID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "user_id": user.ID,
        "name":    user.Name,
        "email":   user.Email,
    })
}