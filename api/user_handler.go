package api

import (
	"hotelreservation/types"

	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	us := types.User{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
	}
	return c.JSON(us)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}