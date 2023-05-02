package user

type User struct {
	FullName    string        `json:"full_name"`
	Department  string        `json:"department"`
	Position    string        `json:"position"`
	ContactInfo []ContactInfo `json:"contact_info"`
}

type Register struct {
	FullName    string        `json:"full_name" validate:"required"`
	Password    string        `json:"password" validate:"required"`
	Department  string        `json:"department"`
	Position    string        `json:"position"`
	ForeignID   string        `json:"foreign_id"`
	OTPEnabled  bool          `json:"otp_enabled"`
	ContactInfo []ContactInfo `json:"contact_info" validate:"dive"`
}

type ContactInfo struct {
	Type  string `json:"type" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}

type Login struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TwoFA struct {
	ID          string `json:"id" validate:"required"`
	OTPPassword string `json:"password" validate:"required"`
}

type OTPData struct {
	UserID    string
	OTPSecret string
}
