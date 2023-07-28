package api

import (
	"hotelreservation/db"
	"hotelreservation/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	room  db.RoomStore
	hotel db.HotelStore
}

func NewHotelHandler(room db.RoomStore, hotel db.HotelStore) *HotelHandler {
	return &HotelHandler{
		room:  room,
		hotel: hotel,
	}
}
func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelId": oid}
	rooms, err := h.room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	hotels, err := h.hotel.GetHotels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	var params types.Hotel
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	hotel, err := h.hotel.InsertHotel(c.Context(), &params)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}