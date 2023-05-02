package postgres

type User struct {
	ID         string
	Department string
	Position   string
	FullName   string
	ForeignID  string

	Password  string
	OTPSecret string

	ContactInfo []ContactInfo `gorm:"->;foreignKey:UserID;references:ID" json:"contact_info"`
}

func (User) TableName() string {
	return "users"
}

type ContactInfo struct {
	UserID string
	Type   string
	Phone  string
}

func (ContactInfo) TableName() string {
	return "contact_info"
}
