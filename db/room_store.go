package db

import (
	"context"
	"hotelreservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
const RoomColl = "rooms"
type RoomStore interface {
	GetRoomByID(context.Context, string) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
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
func (m *MongoRoomStore) GetRoomByID(ctx context.Context, id string) (*types.Room, error) {
	var room types.Room
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	if err := m.coll.FindOne(ctx, filter).Decode(&room); err != nil {
		return nil, err
	}
	return &room, nil
}

func (m *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	var rooms []*types.Room
	cur, err := m.coll.Find(ctx, filter)
	println("cur", cur)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
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