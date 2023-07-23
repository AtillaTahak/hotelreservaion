package db

import (
	"context"
	"hotelreservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
const hotelColl = "hotels"
type HotelStore interface {
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	GetHotels(context.Context) ([]*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	//DeleteHotel(context.Context, string) error
	UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error
	Dropper
}

type MongoHotelStore struct {
	coll   *mongo.Collection
	client *mongo.Client
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(hotelColl),
	}
}

func (m *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cur, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx,&hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (m *MongoHotelStore) Drop(ctx context.Context) error {
	return m.coll.Drop(ctx)
}
func (m *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil
	}
	return nil
	
}

func (m *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	var hotel types.Hotel
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = m.coll.FindOne(ctx, primitive.M{"_id": objID}).Decode(&hotel)
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (m *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := m.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}