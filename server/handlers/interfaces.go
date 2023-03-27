package handlers

import "context"

type Telegram interface {
	SendNotification(ctx context.Context, text string, userID int64) error
}
