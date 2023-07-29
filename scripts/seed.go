package main

import (
	"context"
	"hotelreservation/db"
	"hotelreservation/types"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var (
	client *mongo.Client
	roomStore db.RoomStore
	hotelStore db.HotelStore
	ctx = context.Background()
)
func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:   []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size: 	"small",
			Price: 100,
		},
		{
			Size: 	"normal",
			Price: 200,
		},
		{
			Size: 	"big",
			Price: 300,
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelId = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func main() {
	seedHotel("Hilton", "New York",3)
	seedHotel("Marriot", "New York",5)
	

}

func init(){
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err := client.Database(db.DBNAME).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)

}