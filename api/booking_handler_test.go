package api

import (
	"encoding/json"
	"hotelreservation/api/middleware"
	"hotelreservation/db/fixtures"
	"hotelreservation/types"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

//online get booking Admin

func TestAdminGetBookings(test *testing.T) {
	db := setupUserHandlerTest(test)
	defer db.tearDown(test)

	var (
		adminUser  = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "john", "doe", true)
		hotel          = fixtures.AddHotel(db.Store, "Hilton", "New York", 3, nil)
		room           = fixtures.AddRoom(db.Store, "small", 100, false, hotel.ID)
		from           = time.Now()
		to             = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, to)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Access-Token", CreateTokenFromUser(adminUser))
	if err != nil {
		test.Fatal(err)
	}
	resq, err := app.Test(req)

	if err != nil {
		test.Fatal(err)
	}
	if resq.StatusCode != http.StatusOK {
		test.Fatalf("expected http status of 200 but got %d", resq.StatusCode)
	}
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.Token))
	var bookings []*types.Booking
	if err := json.NewDecoder(resq.Body).Decode(&bookings); err != nil {
		test.Fatal(err)
	}
	if len(bookings) != 1 {
		test.Fatalf("expected 1 booking but got %d", len(bookings))
	}
}
