package slack

import (
	"context"
	"fmt"

	goslack "github.com/slack-go/slack"
)

const (
	notifierType = "slack"
)

type Slack struct {
	webhook string
}

func New(webhook string) *Slack {
	return &Slack{webhook}
}

func (s *Slack) Type() string {
	return notifierType
}

func (s *Slack) Info(ctx context.Context, message string) error {
	msg := messageInfo(message)
	return s.post(ctx, msg)
}

func (s *Slack) Warn(ctx context.Context, message string) error {
	msg := messageWarn(message)
	return s.post(ctx, msg)
}

func (s *Slack) post(ctx context.Context, msg *goslack.WebhookMessage) error {
	if err := goslack.PostWebhookContext(ctx, s.webhook, msg); err != nil {
		return fmt.Errorf("failed to call webhook: %w", err)
	}

	return nil
}

func messageInfo(message string) *goslack.WebhookMessage {
	return &goslack.WebhookMessage{
		Attachments: []goslack.Attachment{
			{
				Color:      "good",
				Fallback:   "fsw-calendar info notification",
				Text:       message,
				MarkdownIn: []string{"text"},
			},
		},
	}
}

func messageWarn(message string) *goslack.WebhookMessage {
	return &goslack.WebhookMessage{
		Attachments: []goslack.Attachment{
			{
				Color:      "warning",
				Fallback:   "fsw-calendar warning notification",
				Text:       message,
				MarkdownIn: []string{"text"},
			},
		},
	}
}
