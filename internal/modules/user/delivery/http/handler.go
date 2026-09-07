package http

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/modules/user/dto"
	"github.com/hurbbiee/todo-list-backend/internal/modules/user/service"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/hurbbiee/todo-list-backend/internal/shared/handler"
	"github.com/hurbbiee/todo-list-backend/internal/shared/helper"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
	"github.com/hurbbiee/todo-list-backend/internal/shared/validator"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	req := new(dto.CreateUserRequest)

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

	err := h.service.Create(
		c.Context(),
		*req,
	)
	if err != nil {
		log.Println(err)
		return handler.HandleError(c, err)
	}

	return response.CreatedJSON(c, "สร้างผู้ใช้งานสำเร็จ")
}

func (h *UserHandler) Profile(c *fiber.Ctx) error {

	actionBy, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	user, err := h.service.Profile(c.Context(), actionBy)
	if err != nil {
		log.Println(err)
		if errors.Is(err, enum.ErrUserNotFound) {
			return response.ErrorJSON(
				c,
				response.BadRequest(enum.UserNotFound),
			)
		}

		return response.ErrorJSON(
			c,
			response.InternalError(enum.Internal),
		)
	}

	return response.SuccessJSON(c, user)
}
