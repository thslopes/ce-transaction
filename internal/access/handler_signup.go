package access

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleSignup(c *fiber.Ctx) error {
	var request struct {
		Token                string `json:"token"`
		Name                 string `json:"name"`
		Email                string `json:"email"`
		Phone                string `json:"phone"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid_request", "Request body must be valid JSON")
	}
	user, err := s.service.RegisterUser(
		request.Token,
		request.Name,
		request.Email,
		request.Phone,
		request.Password,
		request.PasswordConfirmation,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidToken):
			return writeError(c, fiber.StatusUnauthorized, "invalid_token", "Registration token is invalid")
		case errors.Is(err, ErrExpiredToken):
			return writeError(c, fiber.StatusUnauthorized, "expired_token", "Registration token is expired")
		case errors.Is(err, ErrUserAlreadyExists):
			return writeError(c, fiber.StatusConflict, "conflict", "User already exists")
		case errors.Is(err, ErrValidation):
			return writeError(c, fiber.StatusBadRequest, "validation_error", "Request validation failed")
		default:
			return writeError(c, fiber.StatusInternalServerError, "internal_error", "Internal server error")
		}
	}
	return writeJSON(c, fiber.StatusCreated, userResponse(user))
}