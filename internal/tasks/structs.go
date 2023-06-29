package tasks

type Task struct {
	ObjectName          string   `json:"ОbjectName"`
	TaskName            string   `json:"TaskName"`
	UUID                string   `json:"Uid" validate:"required"`
	Date                string   `json:"Date"`
	AuthorID            string   `json:"AuthorId" validate:"required"`
	CreatorID           string   `json:"CreatorId"`
	EndUserID           string   `json:"EndUserId" validate:"required"`
	DeadlineDate        string   `json:"DeadlineDate"`
	TaskInfo            string   `json:"TaskInfo"`
	AgreeStatus         *bool    `json:"AgreeStatus"`
	LinkedTaskID        string   `json:"LinkedTaskId"`
	ApprovalList        []string `json:"ApprovalList"`
	Comment             string   `json:"Comment"`
	EnableDeadDateShift bool     `json:"EnableDeadDateShift"`
	LayoutType          int64    `json:"LayoutType"`
}

type Vacation struct {
	ObjectName          string   `json:"ОbjectName"`
	TaskName            string   `json:"TaskName"`
	UUID                string   `json:"Uid" validate:"required"`
	Date                string   `json:"Date"`
	Author              string   `json:"Author"`
	AuthorID            string   `json:"AuthorId" validate:"required"`
	EndUser             string   `json:"EndUser"`
	EndUserID           string   `json:"EndUserId" validate:"required"`
	DeadlineDate        string   `json:"DeadlineDate"`
	TaskInfo            string   `json:"TaskInfo"`
	AgreeStatus         *bool    `json:"AgreeStatus"`
	LinkedTaskID        string   `json:"LinkedTaskId"`
	ApprovalList        []string `json:"ApprovalList"`
	Comment             string   `json:"Comment"`
	EnableDeadDateShift bool     `json:"EnableDeadDateShift"`
	LayoutType          int64    `json:"LayoutType"`
	PeriodStart         string   `json:"PeriodStart"`
	PeriodEnd           string   `json:"PeriodEnd"`
	HolidayMayker       string   `json:"HolidayMayker"`
	Substitutional      string   `json:"Substitutional"`
}

type Payment struct {
	ObjectName          string   `json:"ОbjectName"`
	TaskName            string   `json:"TaskName"`
	UUID                string   `json:"Uid" validate:"required"`
	Date                string   `json:"Date"`
	Author              string   `json:"Author"`
	AuthorID            string   `json:"AuthorId" validate:"required"`
	EndUser             string   `json:"EndUser"`
	EndUserID           string   `json:"EndUserId" validate:"required"`
	DeadlineDate        string   `json:"DeadlineDate"`
	TaskInfo            string   `json:"TaskInfo"`
	AgreeStatus         *bool    `json:"AgreeStatus"`
	LinkedTaskID        string   `json:"LinkedTaskId"`
	ApprovalList        []string `json:"ApprovalList"`
	Comment             string   `json:"Comment"`
	EnableDeadDateShift bool     `json:"EnableDeadDateShift"`
	LayoutType          int64    `json:"LayoutType"`
	Kontragent          string   `json:"Kontragent"`
	Organization        string   `json:"Organization"`
	Sum                 string   `json:"Sum"`
	PaymentDate         string   `json:"PaymentDate"`
	PaymentPurpose      string   `json:"PaymentPurpose"`
}
