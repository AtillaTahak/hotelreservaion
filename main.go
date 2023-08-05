package main

import (
	"context"
	"flag"
	"hotelreservation/api"
	"hotelreservation/api/middleware"
	"hotelreservation/db"
	"log"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"

)
var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		return ctx.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address for the server")
	flag.Parse()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	//load env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	//handler initialization
	var(
		hotelStore = db.NewMongoHotelStore(client)
		roomStore = db.NewMongoRoomStore(client, hotelStore)
		userStore = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store = &db.Store{
			Room: roomStore,
			Hotel: hotelStore,
			User: userStore,
			Booking: bookingStore,
		}
		userHandler = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		authHandler = api.NewAuthHandler(userStore)
		roomHandler = api.NewRoomHandler(store)
		app = fiber.New(config)
		apiV1 = app.Group("/api/v1",middleware.JWTAuthentication(userStore))
		auth = app.Group("/api")
	

	)
	// auth handlers
	auth.Post("/auth", authHandler.HandleAuthenticate)
	//user handlers
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)

	//hotel handlers
	//apiV1.Delete("/hotel/:id", hotelHandler.HandleDeleteHotel)
	apiV1.Post("/hotel", hotelHandler.HandlePostHotel)
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	//apiV1.Put("/hotel/:id", hotelHandler.HandlePutHotel)

	apiV1.Post("/hotel/:id/book", roomHandler.HandleBookRoom)
	app.Listen(*listenAddr)
}
