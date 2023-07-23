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
func seedHotel(name, location string) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:   []primitive.ObjectID{},
	}
	rooms := []types.Room{
		{
			Type:      types.Single,
			BasePrice: 100,
		},
		{
			Type:      types.Double,
			BasePrice: 200,
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
	seedHotel("Hilton", "New York")
	seedHotel("Marriot", "New York")
	

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