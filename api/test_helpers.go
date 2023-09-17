package api

import (
	"context"
	"hotelreservation/db"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) tearDown(t *testing.T) {
	if err := tdb.client.Database(db.TestDB).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

// TODO
func setupUserHandlerTest(t *testing.T) *testdb {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(db.TestURI))
	if err != nil {
		t.Fatal(err)
	}
	return &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   db.NewMongoHotelStore(client,db.TestDB),
			User:    db.NewMongoUserStore(client,db.TestDB),
			Room:    db.NewMongoRoomStore(client, db.NewMongoHotelStore(client,db.TestDB),db.TestDB),
			Booking: db.NewMongoBookingStore(client,db.TestDB),
		},
	}
}
