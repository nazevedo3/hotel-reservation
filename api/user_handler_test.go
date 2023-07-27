package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nazevedo3/hotel-reservation/types"
)

// const (
// 	testmongouri = "mongodb://localhost:27017"
// 	dname        = "hotel-reservation-test"
// )

// type testdb struct {
// 	db.UserStore
// }

// func (tdb *testdb) teardown(t *testing.T) {
// 	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
// 		t.Fatal(err)
// 	}

// }

// func setup(t *testing.T) *testdb {
// 	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testmongouri))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return &testdb{
// 		UserStore: db.NewMongoUserStore(client),
// 	}
// }

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "James",
		LastName:  "Bond",
		Password:  "kjfljdifjid",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the encrypted password to not be included in JSON response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname to be %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname to be %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email to be %s but got %s", params.Email, user.Email)
	}

}
