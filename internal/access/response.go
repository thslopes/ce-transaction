package access

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func userResponse(user *User) map[string]any {
	response := map[string]any{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"phone":    user.Phone,
		"status":   user.Status,
		"profiles": user.Profiles,
	}
	if !user.CreatedAt.IsZero() {
		response["created_at"] = user.CreatedAt.Format(time.RFC3339)
	}
	if user.ValidatedAt != nil {
		response["validated_at"] = user.ValidatedAt.Format(time.RFC3339)
	}
	if user.ValidatedBy != "" {
		response["validated_by"] = user.ValidatedBy
	}
	return response
}

func writeJSON(c *fiber.Ctx, statusCode int, body any) error {
	return c.Status(statusCode).JSON(body)
}

func writeError(c *fiber.Ctx, statusCode int, code, message string) error {
	return writeJSON(c, statusCode, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
			"details": map[string]any{},
		},
	})
}
