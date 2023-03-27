package structs

type Register struct {
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	OTPEnabled bool   `json:"otp_enabled"`
}

type Login struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type TwoFA struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
