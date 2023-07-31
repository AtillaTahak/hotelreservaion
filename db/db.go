package db

const (
	DBNAME = "hotel-reservation"
	DBURI  = "mongodb://localhost:27017"
	TestDB = "test"
	TestURI = "mongodb://localhost:27017"
)

type Store struct {
	User  UserStore
	Room  RoomStore
	Hotel HotelStore
}
