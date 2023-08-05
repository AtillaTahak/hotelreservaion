package db

import (
	"context"
	"hotelreservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const BookingColl = "bookings"
type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	//DeleteBooking(context.Context, string) error
	//UpdateBooking(ctx context.Context, filter bson.M, params types.UpdateBookingParams) error
	Dropper
}


type MongoBookingStore struct {
	BookingStore
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(BookingColl),
	}
}

func (m *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := m.coll.InsertOne(ctx, booking);
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (m *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	var bookings []*types.Booking
	cur, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}