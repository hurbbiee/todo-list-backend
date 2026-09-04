package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/modules/auth/dto"
	"github.com/hurbbiee/todo-list-backend/internal/modules/auth/service"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
	"github.com/hurbbiee/todo-list-backend/internal/shared/validator"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(dto.AuthRequest)

	// 1. parse request body
	if err := c.BodyParser(req); err != nil {
		return response.ErrorJSON(
			c,
			response.BadRequest(enum.BadRequest),
		)
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		return response.ErrorJSON(
			c,
			response.ValidationError(enum.ValidationError, errs),
		)
	}

	token, err := h.service.Login(
		c.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	return response.SuccessJSON(c, fiber.Map{
		"access_token": token,
	})
}
