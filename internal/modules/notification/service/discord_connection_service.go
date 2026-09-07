package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/dto"
	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/model"
	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/repository"
)

const defaultRemindBeforeMinutes = 60

var (
	ErrInvalidDiscordSettings    = errors.New("invalid Discord settings")
	ErrInvalidDiscordWebhookURL  = errors.New("invalid Discord webhook URL")
	ErrDiscordConnectionNotFound = errors.New("Discord connection not found")
)

type SecretCipher interface {
	Encrypt(plaintext string) ([]byte, error)
	Decrypt(ciphertext []byte) (string, error)
}

type DiscordConnectionService struct {
	repo   repository.DiscordConnectionRepository
	cipher SecretCipher
}

func NewDiscordConnectionService(
	repo repository.DiscordConnectionRepository,
	cipher SecretCipher,
) *DiscordConnectionService {
	return &DiscordConnectionService{
		repo:   repo,
		cipher: cipher,
	}
}

func (s *DiscordConnectionService) Save(
	ctx context.Context,
	req dto.UpsertDiscordConnectionRequest,
	userID int64,
) error {
	if req.IsEnabled == nil ||
		req.NotifyTodoCreated == nil ||
		req.NotifyTodoCompleted == nil ||
		req.NotifyBeforeDue == nil ||
		req.RemindBeforeMinutes < 1 ||
		req.RemindBeforeMinutes > 10080 {
		return ErrInvalidDiscordSettings
	}

	if !isDiscordWebhookURL(req.WebhookURL) {
		return ErrInvalidDiscordWebhookURL
	}

	encryptedWebhookURL, err := s.cipher.Encrypt(req.WebhookURL)
	if err != nil {
		return fmt.Errorf("encrypt Discord webhook URL: %w", err)
	}

	connection := model.DiscordConnection{
		UserID:              userID,
		WebhookURLEncrypted: encryptedWebhookURL,
		IsEnabled:           *req.IsEnabled,
		NotifyTodoCreated:   *req.NotifyTodoCreated,
		NotifyTodoCompleted: *req.NotifyTodoCompleted,
		NotifyBeforeDue:     *req.NotifyBeforeDue,
		RemindBeforeMinutes: req.RemindBeforeMinutes,
	}

	return s.repo.Upsert(ctx, connection)
}

func (s *DiscordConnectionService) Status(
	ctx context.Context,
	userID int64,
) (dto.DiscordConnectionResponse, error) {
	connection, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return dto.DiscordConnectionResponse{}, err
	}

	if connection == nil {
		return dto.DiscordConnectionResponse{
			Connected:           false,
			IsEnabled:           true,
			NotifyTodoCreated:   true,
			NotifyTodoCompleted: true,
			NotifyBeforeDue:     true,
			RemindBeforeMinutes: defaultRemindBeforeMinutes,
		}, nil
	}

	return dto.DiscordConnectionResponse{
		Connected:           true,
		IsEnabled:           connection.IsEnabled,
		NotifyTodoCreated:   connection.NotifyTodoCreated,
		NotifyTodoCompleted: connection.NotifyTodoCompleted,
		NotifyBeforeDue:     connection.NotifyBeforeDue,
		RemindBeforeMinutes: connection.RemindBeforeMinutes,
	}, nil
}

func (s *DiscordConnectionService) UpdateEnabled(
	ctx context.Context,
	req dto.UpdateDiscordEnabledRequest,
	userID int64,
) error {
	if req.IsEnabled == nil {
		return ErrInvalidDiscordSettings
	}

	found, err := s.repo.UpdateEnabled(ctx, userID, *req.IsEnabled)
	if err != nil {
		return err
	}
	if !found {
		return ErrDiscordConnectionNotFound
	}

	return nil
}

func (s *DiscordConnectionService) Disconnect(
	ctx context.Context,
	userID int64,
) error {
	_, err := s.repo.DeleteByUserID(ctx, userID)
	return err
}

func isDiscordWebhookURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Scheme != "https" || parsedURL.User != nil {
		return false
	}
	if parsedURL.Port() != "" || parsedURL.RawQuery != "" || parsedURL.Fragment != "" {
		return false
	}

	host := strings.ToLower(parsedURL.Hostname())
	allowedHost := host == "discord.com" ||
		host == "ptb.discord.com" ||
		host == "canary.discord.com" ||
		host == "discordapp.com"
	if !allowedHost {
		return false
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	return len(parts) == 4 &&
		parts[0] == "api" &&
		parts[1] == "webhooks" &&
		parts[2] != "" &&
		parts[3] != ""
}
