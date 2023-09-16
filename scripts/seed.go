package main

import (
	"context"
	"fmt"
	"hotelreservation/api"
	"hotelreservation/db"
	"hotelreservation/db/fixtures"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func main() {
	var err error
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
		Hotel:   db.NewMongoHotelStore(client),
	}
	user := fixtures.AddUser(store, "james", "foo", false)
	admin := fixtures.AddUser(store, "admin", "admin", true)
	hotel := fixtures.AddHotel(store, "Hilton", "New York", 3, nil)
	room := fixtures.AddRoom(store, "small", 100, false, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ->",booking)
	fmt.Println("user ->",api.CreateTokenFromUser(user))
	fmt.Println("admin ->",api.CreateTokenFromUser(admin))
	
}
