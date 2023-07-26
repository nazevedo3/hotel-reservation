package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nazevedo3/hotel-reservation/db"
	"github.com/nazevedo3/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	ToDate     time.Time `json:"toDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.ToDate) {
		return fmt.Errorf("cannot book a room in the past")
	}

	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}

	ok, err = h.isRoomAvailableForBooking(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked.", c.Params("id")),
		})
	}

	// fmt.Println(bookings)

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		ToDate:     params.ToDate,
		NumPersons: params.NumPersons,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", booking)
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	filter := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"toDate": bson.M{
			"$lte": params.ToDate,
		},
	}

	bookings, err := h.store.Booking.GetBookings(ctx, filter)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil

}
