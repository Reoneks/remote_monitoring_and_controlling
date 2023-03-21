package telegram

import "context"

type Postgres interface {
	BindTelegramUser(ctx context.Context, userID string, telegramUserIDs int64) error
}
