package api

import (
	"errors"
	"fmt"
	"hotelreservation/db"
	"hotelreservation/types"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}
type AuthParams  struct{
	Email string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	User *types.User `json:"user"`
	Token string `json:"token"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error{
	var AuthParams AuthParams
	if err := c.BodyParser(&AuthParams); err != nil {
		return err
	}
	user , err := h.userStore.GetUserByEmail(c.Context(),AuthParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Credentials"})
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword, AuthParams.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Credentials"})
	}
	resp := AuthResponse{
		User: user,
		Token: createTokenFromUser(user),
	}
	return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expiry := now.Add(time.Hour * 4)
	claims := jwt.MapClaims{
		"id": user.ID,
		"email":  user.Email,
		"expiry": expiry.Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	token, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		fmt.Println(err)
	}
	return token

}