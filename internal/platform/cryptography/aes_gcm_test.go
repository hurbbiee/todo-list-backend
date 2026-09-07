package cryptography

import (
	"bytes"
	"errors"
	"testing"
)

func TestAESGCMEncryptDecrypt(t *testing.T) {
	key := bytes.Repeat([]byte{1}, AES256KeySize)
	cipher, err := NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM() error = %v", err)
	}

	const webhookURL = "https://discord.com/api/webhooks/123/token"
	ciphertext, err := cipher.Encrypt(webhookURL)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	if bytes.Contains(ciphertext, []byte(webhookURL)) {
		t.Fatal("ciphertext contains plaintext webhook URL")
	}

	decrypted, err := cipher.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if decrypted != webhookURL {
		t.Fatalf("Decrypt() = %q, want %q", decrypted, webhookURL)
	}
}

func TestAESGCMUsesDifferentNonceForEachEncryption(t *testing.T) {
	key := bytes.Repeat([]byte{2}, AES256KeySize)
	cipher, err := NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM() error = %v", err)
	}

	first, err := cipher.Encrypt("same plaintext")
	if err != nil {
		t.Fatalf("first Encrypt() error = %v", err)
	}
	second, err := cipher.Encrypt("same plaintext")
	if err != nil {
		t.Fatalf("second Encrypt() error = %v", err)
	}

	if bytes.Equal(first, second) {
		t.Fatal("Encrypt() reused ciphertext; expected a unique random nonce")
	}
}

func TestAESGCMRejectsInvalidKeyLength(t *testing.T) {
	_, err := NewAESGCM([]byte("too-short"))
	if err == nil {
		t.Fatal("NewAESGCM() error = nil, want invalid key length error")
	}
}

func TestAESGCMRejectsTamperedCiphertext(t *testing.T) {
	key := bytes.Repeat([]byte{3}, AES256KeySize)
	cipher, err := NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM() error = %v", err)
	}

	ciphertext, err := cipher.Encrypt("secret")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	ciphertext[len(ciphertext)-1] ^= 1

	if _, err := cipher.Decrypt(ciphertext); err == nil {
		t.Fatal("Decrypt() error = nil, want authentication error")
	}
}

func TestAESGCMRejectsEmptyPlaintext(t *testing.T) {
	key := bytes.Repeat([]byte{4}, AES256KeySize)
	cipher, err := NewAESGCM(key)
	if err != nil {
		t.Fatalf("NewAESGCM() error = %v", err)
	}

	_, err = cipher.Encrypt("")
	if !errors.Is(err, ErrEmptyPlaintext) {
		t.Fatalf("Encrypt() error = %v, want %v", err, ErrEmptyPlaintext)
	}
}
