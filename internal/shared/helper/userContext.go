package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
)

func GetUserID(c *fiber.Ctx) (int64, error) {
	value := c.Locals("userId")
	if value == nil {
		return 0, response.Unauthorized(enum.AuthUnauthorized)
	}

	actionBy, ok := value.(int64)
	if !ok {
		return 0, response.InternalError(enum.Internal)
	}

	return actionBy, nil
}
