package router

import "github.com/gofiber/fiber/v3"

func (r *Router) AddTeam(ctx fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
