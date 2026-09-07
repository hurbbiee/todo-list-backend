package http

import (
	"context"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/dto"
	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/service"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/hurbbiee/todo-list-backend/internal/shared/helper"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
	"github.com/hurbbiee/todo-list-backend/internal/shared/validator"
)

type DiscordConnectionService interface {
	Save(
		ctx context.Context,
		req dto.UpsertDiscordConnectionRequest,
		userID int64,
	) error

	Status(
		ctx context.Context,
		userID int64,
	) (dto.DiscordConnectionResponse, error)

	UpdateEnabled(
		ctx context.Context,
		req dto.UpdateDiscordEnabledRequest,
		userID int64,
	) error

	Disconnect(
		ctx context.Context,
		userID int64,
	) error
}

type DiscordConnectionHandler struct {
	service DiscordConnectionService
}

func NewDiscordConnectionHandler(
	service DiscordConnectionService,
) *DiscordConnectionHandler {
	return &DiscordConnectionHandler{service: service}
}

func (h *DiscordConnectionHandler) Save(c *fiber.Ctx) error {
	req := new(dto.UpsertDiscordConnectionRequest)
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

	userID, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	if err := h.service.Save(c.Context(), *req, userID); err != nil {
		return handleDiscordServiceError(c, err)
	}

	return response.SuccessJSON(c, fiber.Map{
		"message": "บันทึกการเชื่อมต่อ Discord สำเร็จ",
	})
}

func (h *DiscordConnectionHandler) Status(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	status, err := h.service.Status(c.Context(), userID)
	if err != nil {
		return handleDiscordServiceError(c, err)
	}

	return response.SuccessJSON(c, status)
}

func (h *DiscordConnectionHandler) UpdateEnabled(c *fiber.Ctx) error {
	req := new(dto.UpdateDiscordEnabledRequest)
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

	userID, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	if err := h.service.UpdateEnabled(c.Context(), *req, userID); err != nil {
		return handleDiscordServiceError(c, err)
	}

	return response.SuccessJSON(c, fiber.Map{
		"message": "อัปเดตสถานะการแจ้งเตือน Discord สำเร็จ",
	})
}

func (h *DiscordConnectionHandler) Disconnect(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	if err := h.service.Disconnect(c.Context(), userID); err != nil {
		return handleDiscordServiceError(c, err)
	}

	return response.SuccessJSON(c, fiber.Map{
		"message": "ยกเลิกการเชื่อมต่อ Discord สำเร็จ",
	})
}

func handleDiscordServiceError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidDiscordSettings):
		return response.ErrorJSON(
			c,
			response.BadRequest(enum.DiscordInvalidSettings),
		)
	case errors.Is(err, service.ErrInvalidDiscordWebhookURL):
		return response.ErrorJSON(
			c,
			response.BadRequest(enum.DiscordInvalidWebhookURL),
		)
	case errors.Is(err, service.ErrDiscordConnectionNotFound):
		return response.ErrorJSON(
			c,
			response.NewAppError(
				fiber.StatusNotFound,
				enum.DiscordConnectionNotFound,
			),
		)
	default:
		log.Printf("Discord connection request failed: %v", err)
		return response.ErrorJSON(
			c,
			response.InternalError(enum.Internal),
		)
	}
}
