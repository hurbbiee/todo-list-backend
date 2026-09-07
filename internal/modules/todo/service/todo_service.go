package service

import (
	"context"
	"log"
	"time"

	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/dto"
	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/event"
	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/repository"
)

type EventPublisher interface {
	PublishJSON(
		ctx context.Context,
		routingKey string,
		payload any,
	) error
}

type TodoService struct {
	repo      repository.TodoRepository
	publisher EventPublisher
}

func NewTodoService(
	repo repository.TodoRepository,
	publisher EventPublisher,
) *TodoService {
	return &TodoService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *TodoService) Search(
	ctx context.Context,
	req dto.SearchTodoRequest,
	actionBy int64,
) ([]dto.SearchTodoResponse, int64, error) {
	return s.repo.Search(ctx, req, actionBy)
}

func (s *TodoService) Create(
	ctx context.Context,
	req dto.CreateTodoRequest,
	actionBy int64,
) error {
	todoID, err := s.repo.Create(
		ctx,
		req,
		actionBy,
	)
	if err != nil {
		return err
	}

	todoCreatedEvent := event.TodoCreated{
		EventType:   event.TodoCreatedRoutingKey,
		TodoID:      todoID,
		UserID:      actionBy,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueAt:       req.DueDate,
		CreatedVia:  req.CreatedVia,
		OccurredAt:  time.Now().UTC(),
	}

	if err := s.publisher.PublishJSON(
		ctx,
		event.TodoCreatedRoutingKey,
		todoCreatedEvent,
	); err != nil {
		// Todo ถูกบันทึกสำเร็จแล้ว
		// จึงไม่ควร return error ให้ frontend เข้าใจว่าสร้างไม่สำเร็จ
		log.Printf(
			"publish todo.created for todoID=%d: %v",
			todoID,
			err,
		)
	}

	return nil
}

func (s *TodoService) CountStatus(
	ctx context.Context,
	actionBy int64,
) (dto.CountTodoResponse, error) {
	return s.repo.CountStatus(ctx, actionBy)
}

func (s *TodoService) Update(
	ctx context.Context,
	req dto.UpdateTodoRequest,
	id int64,
	actionBy int64,
) error {
	return s.repo.Update(ctx, req,id, actionBy)
}

func (s *TodoService) UpdateStatus(
	ctx context.Context,
	req dto.UpdateTodoStatus,
	id int64,
	actionBy int64,
) error {
	return s.repo.UpdateStatus(ctx, req, id,actionBy)
}

func (s *TodoService) Delete(
	ctx context.Context,
	id int64,
	actionBy int64,
) error {
	return s.repo.Delete(ctx, id, actionBy)
}
