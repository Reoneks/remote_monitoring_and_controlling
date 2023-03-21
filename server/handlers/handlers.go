package handlers

type Handler struct {
	postgres Postgres
	telegram Telegram
}

func NewHandler(postgres Postgres, telegram Telegram) *Handler {
	return &Handler{postgres: postgres, telegram: telegram}
}
