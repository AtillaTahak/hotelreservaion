package db

import (
	"context"
	"hotelreservation/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const BookingColl = "bookings"
type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	DeleteBooking(context.Context, string) error
	UpdateBooking(context.Context, string, bson.M) error
	Dropper
}


type MongoBookingStore struct {
	BookingStore
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client,getDbName string) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(getDbName).Collection(BookingColl),
	}
}
func (m *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	updateVal := bson.M{"$set": update}

	resp, err := m.coll.UpdateByID(ctx, objID, updateVal)
	if err != nil {
		return err
	}
	if resp.MatchedCount == 0 {
		return fiber.ErrNotFound
	}
	return nil
}

func (m *MongoBookingStore) DeleteBooking(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if err := m.coll.FindOneAndDelete(ctx, bson.M{"_id": objID}).Err(); err != nil {
		return err
	}
	return nil
}

func (m *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	var booking *types.Booking
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := m.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&booking); err != nil {
		return nil, err
	}
	return booking, nil
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