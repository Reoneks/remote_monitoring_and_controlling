package telegram

import (
	"context"
	"fmt"
	"project/config"
	"project/settings"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

type Telegram struct {
	bot      *tele.Bot
	db       DB
	settings tele.Settings

	frontAddr string
}

func (t *Telegram) SendNotification(ctx context.Context, text string, userID int64) error {
	_, err := t.bot.Send(&tele.User{ID: userID}, text,
		&tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{{
				{Text: "Go to", URL: t.frontAddr + ""}, //TODO: change url
			}},
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to send notification: %w", err)
	}

	return nil
}

func (t *Telegram) Start(ctx context.Context) (err error) {
	t.bot, err = tele.NewBot(t.settings)
	if err != nil {
		err = fmt.Errorf("Failed to create new telegran bot: %w", err)
		return
	}

	t.bot.SetCommands([]tele.Command{
		{Text: "/start", Description: "Start this bot"},
		{Text: "/bind", Description: "Bind chat with bot to database users. Example of usage: '/bind <phone>'"},
	})
	t.prepareHandlers(context.Background())
	go t.bot.Start()
	return
}

func (t *Telegram) prepareHandlers(ctx context.Context) {
	t.bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Plese run the 'bind' command to sync bot and database. Example of command usage: '/bind <phone>'")
	})

	t.bot.Handle("/bind", func(c tele.Context) error {
		phones := c.Args()

		user := c.Sender()
		if user == nil {
			return c.Send("Sender not found. Aborting...")
		}

		if len(phones) <= 0 {
			return c.Send("No phones provided. Aborting...")
		}

		err := t.db.BindTelegramUser(ctx, phones[0], user.ID)
		if err != nil {
			log.Error().Str("function", "prepareHandlers (/bind)").Err(err).Msg("Failed to bind telegram to users")
			return c.Send("Failed to bind telegram to users")
		}

		return c.Send("All users successfully binded")
	})
}

func (t *Telegram) Stop(ctx context.Context) error {
	t.bot.Stop()
	return nil
}

func NewTelegram(cfg *config.Config, db DB) *Telegram {
	return &Telegram{
		db: db,
		settings: tele.Settings{
			Token:  cfg.Token,
			Poller: &tele.LongPoller{Timeout: settings.PollerTimeout},
		},
		frontAddr: cfg.FrontendAddr,
	}
}
