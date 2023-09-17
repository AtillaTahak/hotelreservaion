package api

import (
	"encoding/json"
	"fmt"
	"hotelreservation/api/middleware"
	"hotelreservation/db/fixtures"
	"hotelreservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T){
	db := setupUserHandlerTest(t)
	defer db.tearDown(t)

	var (
		nonAuthUser	   = fixtures.AddUser(db.Store, "James", "Foo", false)
		user           = fixtures.AddUser(db.Store, "john", "doe", false)
		hotel          = fixtures.AddHotel(db.Store, "Hilton", "New York", 3, nil)
		room           = fixtures.AddRoom(db.Store, "small", 100, false, hotel.ID)
		from           = time.Now()
		to             = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, to)
		app            = fiber.New()
		route            = app.Group("/", middleware.JWTAuthentication(db.Store.User))
		bookingHandler = NewBookingHandler(db.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	fmt.Print("booking id ->",booking.ID.Hex())
	req := httptest.NewRequest("GET",fmt.Sprintf("/%s",booking.ID.Hex()), nil)
	req.Header.Set("X-Access-Token", CreateTokenFromUser(user))
	resq, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resq.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resq.StatusCode)
	}

	var bookingRes *types.Booking
	if err := json.NewDecoder(resq.Body).Decode(&bookingRes); err != nil {
		t.Fatal(err)
	}
	if bookingRes.ID.Hex() != booking.ID.Hex() {
		t.Fatalf("expected booking id to be %s but got %s", booking.ID.Hex(), bookingRes.ID.Hex())
	}

	reqTest := httptest.NewRequest("GET", fmt.Sprintf("/%s",booking.ID.Hex()), nil)
	reqTest.Header.Set("X-Access-Token", CreateTokenFromUser(nonAuthUser))
	resq, err = app.Test(reqTest)
	if err != nil {
		t.Fatal(err)
	}
	if resq.StatusCode == http.StatusOK {
		t.Fatalf("expected http status of 403 but got %d", resq.StatusCode)
	}




}

//online get booking Admin

func TestAdminGetBookings(test *testing.T) {
	db := setupUserHandlerTest(test)
	defer db.tearDown(test)

	var (
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "john", "doe", false)
		hotel          = fixtures.AddHotel(db.Store, "Hilton", "New York", 3, nil)
		room           = fixtures.AddRoom(db.Store, "small", 100, false, hotel.ID)
		from           = time.Now()
		to             = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, to)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.Store.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Access-Token", CreateTokenFromUser(adminUser))
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

	reqTest := httptest.NewRequest("GET", "/", nil)
	reqTest.Header.Set("X-Access-Token", CreateTokenFromUser(user))
	resq, err = app.Test(reqTest)
	if err != nil {
		test.Fatal(err)
	}
	if resq.StatusCode == http.StatusOK {
		test.Fatalf("expected http status of 403 but got %d", resq.StatusCode)
	}
}
