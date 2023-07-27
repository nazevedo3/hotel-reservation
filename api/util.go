package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazevedo3/hotel-reservation/types"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, ErrUnAuthorized()
	}
	return user, nil
}
