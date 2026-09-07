package discord

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want POST", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Fatalf("Content-Type = %q", r.Header.Get("Content-Type"))
			}

			var payload struct {
				Content         string `json:"content"`
				AllowedMentions struct {
					Parse []string `json:"parse"`
				} `json:"allowed_mentions"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			if payload.Content != "test notification" {
				t.Fatalf("content = %q", payload.Content)
			}
			if payload.AllowedMentions.Parse == nil ||
				len(payload.AllowedMentions.Parse) != 0 {
				t.Fatalf("allowed mentions = %+v", payload.AllowedMentions.Parse)
			}

			w.WriteHeader(http.StatusNoContent)
		},
	))
	defer server.Close()

	client := NewClient(server.Client())
	if err := client.Send(
		context.Background(),
		server.URL,
		"test notification",
	); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestClientSendEmbed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var payload struct {
				Embeds          []Embed `json:"embeds"`
				AllowedMentions struct {
					Parse []string `json:"parse"`
				} `json:"allowed_mentions"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			if len(payload.Embeds) != 1 {
				t.Fatalf("embeds = %d, want 1", len(payload.Embeds))
			}
			if payload.Embeds[0].Title != "Todo created" ||
				payload.Embeds[0].Color != 0xED4245 {
				t.Fatalf("embed = %+v, want title and color", payload.Embeds[0])
			}
			if payload.AllowedMentions.Parse == nil ||
				len(payload.AllowedMentions.Parse) != 0 {
				t.Fatalf("allowed mentions = %+v", payload.AllowedMentions.Parse)
			}

			w.WriteHeader(http.StatusNoContent)
		},
	))
	defer server.Close()

	err := NewClient(server.Client()).SendEmbed(
		context.Background(),
		server.URL,
		Embed{Title: "Todo created", Color: 0xED4245},
	)
	if err != nil {
		t.Fatalf("SendEmbed() error = %v", err)
	}
}

func TestClientSendMarksRateLimitAsRetryable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Retry-After", "1.5")
			w.WriteHeader(http.StatusTooManyRequests)
		},
	))
	defer server.Close()

	err := NewClient(server.Client()).Send(
		context.Background(),
		server.URL,
		"test notification",
	)
	if err == nil {
		t.Fatal("Send() error = nil, want rate-limit error")
	}
	if !IsRetryable(err) {
		t.Fatal("IsRetryable() = false, want true")
	}

	sendErr, ok := err.(*SendError)
	if !ok {
		t.Fatalf("Send() error type = %T", err)
	}
	if sendErr.RetryAfter != 1500*time.Millisecond {
		t.Fatalf("RetryAfter = %v, want 1.5s", sendErr.RetryAfter)
	}
}

func TestClientSendMarksNotFoundAsPermanent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		},
	))
	defer server.Close()

	err := NewClient(server.Client()).Send(
		context.Background(),
		server.URL,
		"test notification",
	)
	if err == nil {
		t.Fatal("Send() error = nil, want not-found error")
	}
	if IsRetryable(err) {
		t.Fatal("IsRetryable() = true, want false")
	}
}
