package api

import (
	"context"
	"fmt"
	"hotelreservation/db"
	"hotelreservation/types"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
	NumPerson int `json:"numPerson"`
}
func (p *BookRoomParams) Validate() error {
	now := time.Now()
	if p.FromDate.Before(now) {
		return fmt.Errorf("fromDate must be after now")
	}
	if p.FromDate.IsZero() {
		return fmt.Errorf("fromDate is required")
	}
	if p.TillDate.IsZero() {
		return fmt.Errorf("tillDate is required")
	}
	if p.NumPerson <= 0 {
		return fmt.Errorf("numPerson must be greater than 0")
	}
	return nil
}
type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		fmt.Println("error primitive.ObjectIDFromHex")
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Credentials"})
	}
	available, err := h.isRoomAvailableForBooking(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !available {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Room is not available","roomID": c.Params("id")} )
	}
	booking := types.Booking{
		UserID: user.ID,
		RoomID: roomID,
		FromDate: params.FromDate,
		TillDate: params.TillDate,
		NumPerson: params.NumPerson,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func(h *RoomHandler) isRoomAvailableForBooking(c context.Context,id primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomId": id,
		"fromDate": bson.M{"$lte": params.FromDate},
		"tillDate": bson.M{"$gte": params.TillDate},
	}
	bookings , err := h.store.Booking.GetBookings(c, where)
	if err != nil {
		return false,err
	}
	ok := len(bookings) == 0
	return ok,nil

}