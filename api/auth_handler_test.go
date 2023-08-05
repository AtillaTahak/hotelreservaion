package api

import (
	"bytes"
	"context"
	"encoding/json"
	"hotelreservation/db"
	"hotelreservation/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jhon@doe.com",
		Password: "supersecurepassword",
	})
	
	if err != nil {
		t.Fatal(err)
	}
	insertTestUserVar, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	user.ID = insertTestUserVar.ID
	return user
}
func TestAuthenticateFailure(t *testing.T) {
	tdb := setupUserHandlerTest(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.Store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "jhon@doe.com",
		Password: "supersecurepasswordasdasd",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected http status of 401 but got %d", resp.StatusCode)
	}
}
func TestAuthenticateSuccess(t *testing.T) {
	tdb := setupUserHandlerTest(t)
	defer tdb.tearDown(t)
	insertedUser := insertTestUser(t, tdb.Store.User)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.Store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "jhon@doe.com",
		Password: "supersecurepassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}

	insertedUser.EncryptedPassword = ""
	if authResp.User.FirstName != insertedUser.FirstName ||
	authResp.User.LastName != insertedUser.LastName ||
	authResp.User.Email != insertedUser.Email {
	 t.Fatalf("expected the user fields to be present in the auth response")
 	}
}
