package fixtures

import (
	"context"
	"fmt"
	"hotelreservation/db"
	"hotelreservation/types"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, startDate, endDate time.Time) *types.Booking {
	booking := types.Booking{
		UserID:    userID,
		RoomID:    roomID,
		FromDate:  startDate,
		TillDate:  endDate,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), &booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func AddRoom(store *db.Store, size string, price float64,ss bool, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		Size: 	size,
		Price: price,
		Seaside: ss,
		HotelId: hotelID,
	}
	insertedRoom, err := store.Room.InsertRoom(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}
func AddHotel(store *db.Store, name , loc string , rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDs = rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:   roomIDs,
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store,fn,ln string,admin bool) *types.User {

	user, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:    fmt.Sprintf("%s@%s.com",fn,ln),
		Password:  fmt.Sprintf("%s_%s",fn,ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin;
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}