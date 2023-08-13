package util

import (
	"hotelreservation/types"
	"github.com/gofiber/fiber/v2"
)

func GetAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fiber.ErrForbidden
	}
	return user, nil
}
