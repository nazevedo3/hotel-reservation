package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazevedo3/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(hotel)

}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	// var pagination db.Pagination
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	filter := bson.M{"rating": params.Rating}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrResourceNotFound("hotels")
	}
	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(params.Pagination.Page),
	}
	return c.JSON(resp)
}
