package response

import "github.com/gofiber/fiber/v2"

type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ErrorJSON(c *fiber.Ctx, err error) error {
	status, resp := FromError(err)
	return c.Status(status).JSON(resp)
}

func SuccessJSON[T any](c *fiber.Ctx, data T) error {
	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"success": true,
		"data":    data,
	})
}

func CreatedJSON(
	c *fiber.Ctx,
	message string,
) error {
	return c.Status(fiber.StatusCreated).JSON(MessageResponse{
		Success: true,
		Message: message,
	})
}
