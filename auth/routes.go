package auth

import (
	"hrms/common"
	"hrms/users"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(r fiber.Router, userRepo users.UserRepository, log *slog.Logger, jwtSecret string) {

	authService := NewAuthService(userRepo, log, jwtSecret)

	r.Post("/register", func(c *fiber.Ctx) error {
		user := new(UserAuth)
		if err := common.ValidateRequest(c, user); err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, err)
		}
		if err := authService.Register(user); err != nil {
			switch err {
			case ErrUserExists:
				return common.ErrorResponse(c, fiber.StatusConflict, err)
			case common.ErrDatabase, ErrSaveUser, ErrHashPassword:
				return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
			default:
				return common.ErrorResponse(c, fiber.StatusBadRequest, err)
			}
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User signed up successfully"})
	})

	r.Post("/login", func(c *fiber.Ctx) error {
		user := new(UserAuth)
		if err := common.ValidateRequest(c, user); err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, err)
		}
		token, err := authService.Login(user)
		if err != nil {
			switch err {
			case ErrGenerateToken:
				return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
			default:
				return common.ErrorResponse(c, fiber.StatusBadRequest, err)
			}
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
	})
}
