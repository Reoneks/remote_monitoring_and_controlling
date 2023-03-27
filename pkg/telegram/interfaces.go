package telegram

import "context"

type DB interface {
	BindTelegramUser(ctx context.Context, userPhone string, telegramUserIDs int64) error
}
