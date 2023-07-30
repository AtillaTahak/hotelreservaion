package api

import (
	"errors"
	"hotelreservation/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(AuthParams.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Credentials"})
	}
	return c.JSON(fiber.Map{"message": "Authenticated"})

}