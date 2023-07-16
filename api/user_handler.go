package api

import (
	"context"
	"hotelreservation/db"
	"hotelreservation/types"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	us := types.User{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
	}
	return c.JSON(us)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
		ctx = context.Background()
	)
	user, err := h.userStore.GetUserByID(ctx,id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(user)
}