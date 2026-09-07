package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	DefaultTimeout       = 10 * time.Second
	maxRateLimitBodySize = 4 << 10
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPDoer
}

type Embed struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Color       int          `json:"color,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Footer      *EmbedFooter `json:"footer,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedFooter struct {
	Text string `json:"text"`
}

type webhookPayload struct {
	Content         string  `json:"content,omitempty"`
	Embeds          []Embed `json:"embeds,omitempty"`
	AllowedMentions struct {
		Parse []string `json:"parse"`
	} `json:"allowed_mentions"`
}

type SendError struct {
	StatusCode int
	Retryable  bool
	RetryAfter time.Duration
	Err        error
}

func (e *SendError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf(
			"Discord webhook returned HTTP status %d",
			e.StatusCode,
		)
	}

	return "Discord webhook request failed"
}

func (e *SendError) Unwrap() error {
	return e.Err
}

func IsRetryable(err error) bool {
	var sendErr *SendError
	return errors.As(err, &sendErr) && sendErr.Retryable
}

func NewClient(httpClient HTTPDoer) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: DefaultTimeout}
	}

	return &Client{httpClient: httpClient}
}

func (c *Client) Send(
	ctx context.Context,
	webhookURL string,
	content string,
) error {
	payload := webhookPayload{Content: content}
	return c.send(ctx, webhookURL, payload)
}

func (c *Client) SendEmbed(
	ctx context.Context,
	webhookURL string,
	embed Embed,
) error {
	payload := webhookPayload{Embeds: []Embed{embed}}
	return c.send(ctx, webhookURL, payload)
}

func (c *Client) send(
	ctx context.Context,
	webhookURL string,
	payload webhookPayload,
) error {
	payload.AllowedMentions.Parse = []string{}

	body, err := json.Marshal(payload)
	if err != nil {
		return &SendError{Err: fmt.Errorf("encode Discord payload: %w", err)}
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		webhookURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return &SendError{Err: fmt.Errorf("create Discord request: %w", err)}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "todo-workspace/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &SendError{
			Retryable: true,
			Err:       err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK &&
		resp.StatusCode < http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	retryable := resp.StatusCode == http.StatusRequestTimeout ||
		resp.StatusCode == http.StatusTooManyRequests ||
		resp.StatusCode >= http.StatusInternalServerError

	return &SendError{
		StatusCode: resp.StatusCode,
		Retryable:  retryable,
		RetryAfter: readRetryAfter(resp),
	}
}

func readRetryAfter(resp *http.Response) time.Duration {
	if resp.StatusCode != http.StatusTooManyRequests {
		return 0
	}
	limitedBody := io.LimitReader(resp.Body, maxRateLimitBodySize)
	bodyBytes, _ := io.ReadAll(limitedBody)

	if seconds, err := strconv.ParseFloat(
		resp.Header.Get("Retry-After"),
		64,
	); err == nil && seconds > 0 {
		return time.Duration(seconds * float64(time.Second))
	}

	var body struct {
		RetryAfter float64 `json:"retry_after"`
	}
	if err := json.Unmarshal(bodyBytes, &body); err != nil ||
		body.RetryAfter <= 0 {
		return 0
	}

	return time.Duration(body.RetryAfter * float64(time.Second))
}
