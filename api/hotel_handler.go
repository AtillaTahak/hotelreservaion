package api

import (
	"hotelreservation/db"
	"hotelreservation/types"

	"github.com/gofiber/fiber/v2"
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
type HotelQueryParams struct {
	Rooms bool
	Rating int
}


func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	//TODO: implement query params
	// TODO: implement pagination

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