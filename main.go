package main

import (
	"context"
	"flag"
	"hotelreservation/api"
	"hotelreservation/db"
	"log"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	if err := client.Database(db.DBNAME).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	//handler initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))

	app := fiber.New(config)
	apiV1 := app.Group("/api/v1")

	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	app.Listen(*listenAddr)
}
