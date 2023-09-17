package api

import (
	"bytes"
	"encoding/json"
	"hotelreservation/db/fixtures"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateFailure(t *testing.T) {
	tdb := setupUserHandlerTest(t)
	defer tdb.tearDown(t)
	fixtures.AddUser(tdb.Store,"John","Doe",false)

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
	//insertedUser := insertTestUser(t, tdb.Store.User)
	insertedUser := fixtures.AddUser(tdb.Store,"john","doe",false)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.Store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "john@doe.com",
		Password: "john_doe",
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
