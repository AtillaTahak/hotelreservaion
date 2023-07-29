package db

const (
	DBNAME = "hotel-reservation"
	DBURI  = "mongodb://localhost:27017"
	TestDB = "test"
)

type Store struct {
	User  UserStore
	Room  RoomStore
	Hotel HotelStore
}
