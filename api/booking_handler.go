package api

import (
	"hotelreservation/db"
	"hotelreservation/util"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	user, err := util.GetAuthUser(c)
	if err != nil{
		return err
	}
	if booking.UserID != user.ID {
		return fiber.ErrForbidden
	}
	err = h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true})
	if err != nil {
		return err
	}
	return c.JSON(booking)
}


//TODO: this needs to be admin auth	
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)	
}

//TODO: this needs to be user auth
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	user, err := util.GetAuthUser(c)
	if err != nil{
		return err
	}
	if booking.UserID != user.ID {
		return fiber.ErrForbidden
	}
	return c.JSON(booking)
}