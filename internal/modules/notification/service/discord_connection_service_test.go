package service

import (
	"bytes"
	"context"
	"testing"

	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/dto"
	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/model"
)

type fakeDiscordConnectionRepository struct {
	connection *model.DiscordConnection
	upserted   *model.DiscordConnection
}

func (f *fakeDiscordConnectionRepository) FindByUserID(
	_ context.Context,
	_ int64,
) (*model.DiscordConnection, error) {
	return f.connection, nil
}

func (f *fakeDiscordConnectionRepository) Upsert(
	_ context.Context,
	connection model.DiscordConnection,
) error {
	f.upserted = &connection
	return nil
}

func (f *fakeDiscordConnectionRepository) UpdateEnabled(
	_ context.Context,
	_ int64,
	_ bool,
) (bool, error) {
	return true, nil
}

func (f *fakeDiscordConnectionRepository) DeleteByUserID(
	_ context.Context,
	_ int64,
) (bool, error) {
	return true, nil
}

type fakeSecretCipher struct {
	ciphertext []byte
}

func (f *fakeSecretCipher) Encrypt(_ string) ([]byte, error) {
	return f.ciphertext, nil
}

func (f *fakeSecretCipher) Decrypt(_ []byte) (string, error) {
	return "", nil
}

func boolPointer(value bool) *bool {
	return &value
}

func TestDiscordConnectionServiceSaveStoresEncryptedWebhook(t *testing.T) {
	repo := &fakeDiscordConnectionRepository{}
	cipher := &fakeSecretCipher{ciphertext: []byte("encrypted-webhook")}
	service := NewDiscordConnectionService(repo, cipher)

	req := dto.UpsertDiscordConnectionRequest{
		WebhookURL:          "https://discord.com/api/webhooks/123/token",
		IsEnabled:           boolPointer(true),
		NotifyTodoCreated:   boolPointer(true),
		NotifyTodoCompleted: boolPointer(true),
		NotifyBeforeDue:     boolPointer(true),
		RemindBeforeMinutes: 60,
	}

	if err := service.Save(context.Background(), req, 7); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if repo.upserted == nil {
		t.Fatal("Save() did not call repository Upsert()")
	}
	if repo.upserted.UserID != 7 {
		t.Fatalf("Upsert() user ID = %d, want 7", repo.upserted.UserID)
	}
	if !bytes.Equal(repo.upserted.WebhookURLEncrypted, cipher.ciphertext) {
		t.Fatal("Upsert() did not receive encrypted webhook URL")
	}
	if bytes.Equal(repo.upserted.WebhookURLEncrypted, []byte(req.WebhookURL)) {
		t.Fatal("Upsert() received plaintext webhook URL")
	}
}

func TestDiscordConnectionServiceRejectsNonDiscordWebhook(t *testing.T) {
	repo := &fakeDiscordConnectionRepository{}
	cipher := &fakeSecretCipher{ciphertext: []byte("encrypted-webhook")}
	service := NewDiscordConnectionService(repo, cipher)

	req := dto.UpsertDiscordConnectionRequest{
		WebhookURL:          "https://example.com/api/webhooks/123/token",
		IsEnabled:           boolPointer(true),
		NotifyTodoCreated:   boolPointer(true),
		NotifyTodoCompleted: boolPointer(true),
		NotifyBeforeDue:     boolPointer(true),
		RemindBeforeMinutes: 60,
	}

	if err := service.Save(context.Background(), req, 7); err == nil {
		t.Fatal("Save() error = nil, want invalid Discord webhook error")
	}
	if repo.upserted != nil {
		t.Fatal("Save() called repository for an invalid webhook URL")
	}
}

func TestDiscordConnectionServiceStatusDoesNotExposeWebhook(t *testing.T) {
	repo := &fakeDiscordConnectionRepository{
		connection: &model.DiscordConnection{
			UserID:              7,
			WebhookURLEncrypted: []byte("encrypted-webhook"),
			IsEnabled:           true,
			NotifyTodoCreated:   true,
			NotifyTodoCompleted: false,
			NotifyBeforeDue:     true,
			RemindBeforeMinutes: 60,
		},
	}
	service := NewDiscordConnectionService(repo, &fakeSecretCipher{})

	status, err := service.Status(context.Background(), 7)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if !status.Connected || status.NotifyTodoCompleted {
		t.Fatalf("Status() = %+v, want connected settings without webhook", status)
	}
}

func TestDiscordConnectionServiceStatusReturnsDefaultsWhenDisconnected(t *testing.T) {
	service := NewDiscordConnectionService(
		&fakeDiscordConnectionRepository{},
		&fakeSecretCipher{},
	)

	status, err := service.Status(context.Background(), 7)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if status.Connected {
		t.Fatal("Status().Connected = true, want false")
	}
	if !status.IsEnabled ||
		!status.NotifyTodoCreated ||
		!status.NotifyTodoCompleted ||
		!status.NotifyBeforeDue ||
		status.RemindBeforeMinutes != defaultRemindBeforeMinutes {
		t.Fatalf("Status() defaults = %+v", status)
	}
}
