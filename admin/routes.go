package admin

import (
	"hrms/common"
	"hrms/models"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(r fiber.Router, adminRepo AdminRepository) {
	r.Get("/candidates", func(c *fiber.Ctx) error {
		status := c.Query("status", string(models.StatusPending))
		candidates, err := adminRepo.GetCandidates(models.CandidateStatus(status))
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Status(fiber.StatusOK).JSON(candidates)
	})

	r.Put("/candidates/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")
		candidate := new(models.Candidate)
		if err := common.ValidateRequest(c, candidate); err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, err)
		}
		candidate.ID = id
		updatedCandidate, err := adminRepo.UpdateCandidate(candidate)
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Status(fiber.StatusOK).JSON(updatedCandidate)
	})

	r.Get("/candidates/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")
		candidate, err := adminRepo.GetCandidateByID(id)
		if err != nil {
			if err == ErrCandidateNotFound {
				return common.ErrorResponse(c, fiber.StatusNotFound, err)
			}
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Status(fiber.StatusOK).JSON(candidate)
	})
}
