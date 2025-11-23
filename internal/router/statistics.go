package router

import "github.com/gofiber/fiber/v3"

func (r *Router) GetStatistic(c fiber.Ctx) error {
	stats, err := r.service.GetStatistic(c.Context())
	if err != nil {
		r.logger.WithError(err).Error("failed to get statistics")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get statistics"})
	}
	return c.Status(fiber.StatusOK).JSON(stats)
}
