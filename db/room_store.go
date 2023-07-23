package db

import (
	"context"
	"hotelreservation/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
const RoomColl = "rooms"
type RoomStore interface {
	//GetRoomByID(context.Context, string) (*types.Room, error)
	//GetRooms(context.Context) ([]*types.Room, error)
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	//DeleteRoom(context.Context, string) error
	//UpdateRoom(ctx context.Context, filter bson.M, params types.UpdateRoomParams) error
	Dropper
}

type MongoRoomStore struct {
	coll   *mongo.Collection
	client *mongo.Client

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(RoomColl),
		HotelStore: hotelStore,
	}
}

func (m *MongoRoomStore) Drop(ctx context.Context) error {
	return m.coll.Drop(ctx)
}

func (m *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := m.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)
	filter := primitive.M{"_id": room.HotelId}
	update := primitive.M{"$push": primitive.M{"rooms": room.ID}}
	if err := m.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil

}