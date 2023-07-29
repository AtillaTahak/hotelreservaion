package types

import "go.mongodb.org/mongo-driver/bson/primitive"
type CreateHotelParams struct {
	Name string `json:"name"`
	Location string `json:"location"`
}

type Hotel struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Location string `bson:"location" json:"location"`
	Rooms []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating int `bson:"rating" json:"rating"`
}
type RoomType int

const (
	_ RoomType = iota
	Single
	Double
	Twin
	Suite
)

type Room struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	//Type RoomType `bson:"type" json:"type"`	
	Size string `bson:"size" json:"size"`
	Seaside bool `bson:"seaside" json:"seaside"`
	Price float64 `bson:"price" json:"price"`
	HotelId primitive.ObjectID `bson:"hotelId" json:"hotelId"`
}