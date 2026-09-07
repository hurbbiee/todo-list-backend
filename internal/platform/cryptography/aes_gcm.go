package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

const AES256KeySize = 32

var (
	ErrEmptyPlaintext    = errors.New("plaintext is empty")
	ErrInvalidCiphertext = errors.New("ciphertext is invalid")
)

type AESGCM struct {
	aead cipher.AEAD
}

func NewAESGCM(key []byte) (*AESGCM, error) {
	if len(key) != AES256KeySize {
		return nil, fmt.Errorf(
			"AES-256 key must be %d bytes",
			AES256KeySize,
		)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create AES-GCM: %w", err)
	}

	return &AESGCM{aead: aead}, nil
}

func (c *AESGCM) Encrypt(plaintext string) ([]byte, error) {
	if plaintext == "" {
		return nil, ErrEmptyPlaintext
	}

	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate AES-GCM nonce: %w", err)
	}

	// เก็บ nonce นำหน้า ciphertext เพื่อให้ถอดรหัสได้ภายหลัง
	return c.aead.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

func (c *AESGCM) Decrypt(ciphertext []byte) (string, error) {
	nonceSize := c.aead.NonceSize()
	if len(ciphertext) <= nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce := ciphertext[:nonceSize]
	encryptedData := ciphertext[nonceSize:]

	plaintext, err := c.aead.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt AES-GCM ciphertext: %w", err)
	}

	return string(plaintext), nil
}
