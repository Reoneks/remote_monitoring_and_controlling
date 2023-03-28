package structs

type User struct {
	ID       string
	Phone    string
	Password string

	OTPEnabled bool
	OTPSecret  string

	TelegramUserID int64
}
