package postgres

import "github.com/lib/pq"

type User struct {
	ID         string
	Department string
	Position   string
	FullName   string

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

type TaskType uint8

const (
	TaskT TaskType = iota + 1
	VacationT
	PaymentT
)

type Task struct {
	ObjectName          string
	TaskName            string
	TaskType            TaskType
	UUID                string
	Date                string
	Author              string
	AuthorID            string
	CreatorID           string
	EndUser             string
	EndUserID           string
	DeadlineDate        string
	TaskInfo            string
	AgreeStatus         *bool
	LinkedTaskID        string
	ApprovalList        pq.StringArray `gorm:"type:text[]"`
	Comment             string
	EnableDeadDateShift bool
	LayoutType          int64
	PeriodStart         string
	PeriodEnd           string
	HolidayMayker       string
	Substitutional      string
	Kontragent          string
	Organization        string
	Sum                 string
	PaymentDate         string
	PaymentPurpose      string
}

func (Task) TableName() string {
	return "tasks"
}
