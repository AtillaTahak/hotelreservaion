package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hotelreservation/db"
	"hotelreservation/types"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
const (
	testdburi = "mongodb://localhost:27017"
	testdbname = "hotel-reservation-test"
	)
type testdb struct {
	db.UserStore
}
func (tdb *testdb) tearDown() {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

func setupUserHandlerTest(t *testing.T) *testdb{
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client,testdbname),
	}
}
func TestPostUser(t *testing.T) {
	tdb := setupUserHandlerTest(t)
	
	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/user", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jhondoe@hotmail.com",
		Password: "12345678",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/user",bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bb))

	defer tdb.tearDown()
}