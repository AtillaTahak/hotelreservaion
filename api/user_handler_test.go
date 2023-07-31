package api

import (
	"bytes"
	"context"
	"encoding/json"
	"hotelreservation/db"
	"hotelreservation/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) tearDown(t *testing.T) {
	if err := tdb.client.Database(db.TestDB).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setupUserHandlerTest(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.TestURI))
	if err != nil {
		t.Fatal(err)
	}
	return &testdb{
		client: client,
		Store: &db.Store{
			User: db.NewMongoUserStore(client),
		},
	}
}
func TestPostUser(t *testing.T) {
	tdb := setupUserHandlerTest(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store.User)
	app.Post("/user", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jhondoe@hotmail.com",
		Password:  "12345678",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/user", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if user.FirstName != params.FirstName {
		t.Errorf("Expected %s, got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("Expected %s, got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("Expected %s, got %s", params.Email, user.Email)
	}

	defer tdb.tearDown(t)
}

func TestGetUsers(t *testing.T) {
	tdb := setupUserHandlerTest(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store.User)
	app.Get("/users", userHandler.HandleGetUsers)

	req := httptest.NewRequest("Get", "/user", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var users []types.User
	json.NewDecoder(resp.Body).Decode(&users)
	if len(users) != 0 {
		t.Errorf("Expected %d, got %d", 0, len(users))
	}
	defer tdb.tearDown(t)
}
