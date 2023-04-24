package structs

type Register struct {
	Phone      string `json:"phone" validate:"required"`
	Password   string `json:"password" validate:"required"`
	OTPEnabled bool   `json:"otp_enabled"`
}

type Login struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TwoFA struct {
	Phone       string `json:"phone" validate:"required"`
	OTPPassword string `json:"password" validate:"required"`
}
