package config

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func setRequiredTestEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	t.Setenv("JWT_SECRET", "test-jwt-secret")
	t.Setenv("APP_ENV", "test")
}

func TestLoadDecodesWebhookEncryptionKey(t *testing.T) {
	setRequiredTestEnvironment(t)
	wantKey := bytes.Repeat([]byte{7}, webhookEncryptionKeySize)
	t.Setenv(
		"WEBHOOK_ENCRYPTION_KEY",
		base64.StdEncoding.EncodeToString(wantKey),
	)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !bytes.Equal(cfg.Security.WebhookEncryptionKey, wantKey) {
		t.Fatal("Load() returned an unexpected webhook encryption key")
	}
}

func TestLoadRejectsMissingWebhookEncryptionKey(t *testing.T) {
	setRequiredTestEnvironment(t)
	t.Setenv("WEBHOOK_ENCRYPTION_KEY", "")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "WEBHOOK_ENCRYPTION_KEY is required") {
		t.Fatalf("Load() error = %v, want missing key error", err)
	}
}

func TestLoadRejectsInvalidBase64WebhookEncryptionKey(t *testing.T) {
	setRequiredTestEnvironment(t)
	t.Setenv("WEBHOOK_ENCRYPTION_KEY", "not-base64")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "valid base64") {
		t.Fatalf("Load() error = %v, want base64 error", err)
	}
}

func TestLoadRejectsWrongWebhookEncryptionKeySize(t *testing.T) {
	setRequiredTestEnvironment(t)
	t.Setenv(
		"WEBHOOK_ENCRYPTION_KEY",
		base64.StdEncoding.EncodeToString([]byte("too-short")),
	)

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "decode to 32 bytes") {
		t.Fatalf("Load() error = %v, want key size error", err)
	}
}
