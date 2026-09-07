package http

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/dto"
	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/service"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/hurbbiee/todo-list-backend/internal/shared/handler"
	"github.com/hurbbiee/todo-list-backend/internal/shared/helper"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
	"github.com/hurbbiee/todo-list-backend/internal/shared/validator"
)

type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(s *service.TodoService) *TodoHandler {
	return &TodoHandler{service: s}
}

func (h *TodoHandler) Search(c *fiber.Ctx) error {
	req := new(dto.SearchTodoRequest)

	if err := c.BodyParser(req); err != nil {
		return response.ErrorJSON(
			c,
			response.BadRequest(enum.BadRequest),
		)
	}

	req.ApplyDefault()

	if errs := validator.ValidateStruct(req); errs != nil {
		return response.ErrorJSON(
			c,
			response.ValidationError(enum.ValidationError, errs),
		)
	}

	actionBy, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	todos, total, err := h.service.Search(c.Context(), *req, actionBy)
	if err != nil {
		log.Println(err)
		return response.ErrorJSON(
			c,
			response.InternalError(enum.Internal),
		)
	}

	return response.SuccessJSON(
		c,
		response.Paginate(
			todos,
			req.Page,
			req.Limit,
			total,
		),
	)
}

func (h *TodoHandler) Create(c *fiber.Ctx) error {
	req := new(dto.CreateTodoRequest)

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

	actionBy, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	err = h.service.Create(c.Context(), *req, actionBy)
	if err != nil {
		log.Println(err)
		return handler.HandleError(c, err)
	}

	return response.CreatedJSON(c, "สร้างงานงานสำเร็จ")
}

func (h *TodoHandler) CountStatus(c *fiber.Ctx) error {
	actionBy, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	count, err := h.service.CountStatus(c.Context(), actionBy)
	if err != nil {
		return handler.HandleError(c, err)
	}
	return response.SuccessJSON(c, count)

}

func (h *TodoHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return response.ErrorJSON(
			c,
			response.BadRequest(enum.BadRequest),
		)
	}

	req := new(dto.UpdateTodoRequest)

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

	actionBy, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	err = h.service.Update(c.Context(), *req, id,actionBy)
	if err != nil {
		log.Println(err)
		return handler.HandleError(c, err)
	}

	return response.CreatedJSON(c, "อัปเดทสิ่งที่ต้องทำสำเร็จ")
}

func (h *TodoHandler) UpdateStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return response.ErrorJSON(
			c,
			response.BadRequest(enum.BadRequest),
		)
	}

	req := new(dto.UpdateTodoStatus)

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

	actionBy, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	err = h.service.UpdateStatus(c.Context(), *req, id,actionBy)
	if err != nil {
		log.Println(err)
		return handler.HandleError(c, err)
	}

	return response.CreatedJSON(c, "ปรับสถานะสิ่งที่ต้องทำเป็น เสร็จแล้ว เรียบร้อย")
}

func (h *TodoHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return response.ErrorJSON(
			c,
			response.BadRequest(enum.BadRequest),
		)
	}

	actionBy, err := helper.GetUserID(c)
	if err != nil {
		return response.ErrorJSON(c, err)
	}

	err = h.service.Delete(c.Context(), id, actionBy)
	if err != nil {
		return handler.HandleError(c, err)
	}

	return response.CreatedJSON(c, "ลบสิ่งที่ต้องทำ สำเร็จ")
}
