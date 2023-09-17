package db

import (
	"context"
	"fmt"
	"hotelreservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
const userColl = "users"
type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	Dropper
}
type MongoUserStore struct {
	coll *mongo.Collection
	client *mongo.Client
}

func NewMongoUserStore(client *mongo.Client,getDbName string) *MongoUserStore {
	return &MongoUserStore{
		client:client ,
		coll:  client.Database(getDbName).Collection(userColl),
	}
}

func (m *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	err := m.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (m *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("dropping user collection")
	return m.coll.Drop(ctx)
}
func (m *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	update := bson.M{"$set": params}
	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil
	}
	return nil
	
}
func (m *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (m *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx,&users); err != nil {
		return nil, err
	}
	return users, nil

/* 	var users []*types.User
	cur, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user types.User
		if err := cur.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil */
}

func (m *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// TODO: Maybe its a good idea to handle if we did not delete anything
	// maybe log it or something
	var user types.User
	if err := m.coll.FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return err
	}
	return nil
}
	
func (m *MongoUserStore) GetUserByID(ctx context.Context ,id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}