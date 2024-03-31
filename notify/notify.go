package notify

import "context"

type Notifier interface {
	Type() string
	Info(ctx context.Context, message string) error
	Warn(ctx context.Context, message string) error
}
