package user

type Register struct {
	FullName    string        `json:"FullName" validate:"required"`
	Password    string        `json:"Password" validate:"required"`
	Department  string        `json:"Department"`
	Position    string        `json:"Position"`
	ForeignID   string        `json:"UID" validate:"required"`
	ContactInfo []ContactInfo `json:"ContactInfo" validate:"dive"`
}

type AddAlternativeNumber struct {
	ForeignID   string        `json:"UID" validate:"required"`
	ContactInfo []ContactInfo `json:"ContactInfo" validate:"dive"`
}

type ContactInfo struct {
	Type  string `json:"Type"`
	Phone string `json:"Telephone" validate:"required"`
}

type Login struct {
	Phone    string `json:"Telephone" validate:"required"`
	Password string `json:"Password" validate:"required"`
}

type TwoFA struct {
	ID          string `json:"ID" validate:"required"`
	OTPPassword string `json:"Password" validate:"required"`
}

type OTPData struct {
	UserID    string
	OTPSecret string
}
