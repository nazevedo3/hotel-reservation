package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazevedo3/hotel-reservation/db/fixtures"
	"github.com/nazevedo3/hotel-reservation/types"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "Jimmy", "watercooler", false)
		user           = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		to             = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, to)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		route          = app.Group("/:id", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)
	_ = app
	route.Get("/", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal(err)
	}
	var bookingResponse *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResponse); err != nil {
		t.Fatal(err)
	}
	if bookingResponse.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResponse.ID)
	}
	if bookingResponse.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResponse.UserID)
	}
	fmt.Println(bookingResponse)

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatal(err)
	}

}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		to             = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, to)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response got %d\n", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking but got %d\n", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
	}

	// test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized but got %d\n", resp.StatusCode)
	}

}
