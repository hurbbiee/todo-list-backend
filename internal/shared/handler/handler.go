package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
)

var errorMap = map[error]*response.AppError{

	enum.ErrEmailAlreadyExists: response.BadRequest(enum.EmailAlreadyExists),
}

func HandleError(c *fiber.Ctx, err error) error {
	if res, ok := errorMap[err]; ok {
		return response.ErrorJSON(c, res)
	}

	for e, res := range errorMap {
		if errors.Is(err, e) {
			return response.ErrorJSON(c, res)
		}
	}

	return response.ErrorJSON(
		c,
		response.InternalError(enum.Internal),
	)
}
