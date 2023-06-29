package tasks

import (
	"context"
	"fmt"
	"remote_monitoring_and_controlling/config"
	"remote_monitoring_and_controlling/pkg/postgres"

	"github.com/go-resty/resty/v2"
)

type Service struct {
	db      DB
	resty   *resty.Client
	basAddr string
}

func (u *Service) GetTasks(ctx context.Context, userID string) ([]Task, error) {
	tasks, err := u.db.GetTasks(ctx, userID, postgres.TaskT)
	if err != nil {
		return nil, err
	}

	tasksResp := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		tasksResp = append(tasksResp, Task{
			ObjectName:          task.ObjectName,
			TaskName:            task.TaskName,
			UUID:                task.UUID,
			Date:                task.Date,
			AuthorID:            task.AuthorID,
			CreatorID:           task.CreatorID,
			EndUserID:           task.EndUserID,
			DeadlineDate:        task.DeadlineDate,
			TaskInfo:            task.TaskInfo,
			AgreeStatus:         task.AgreeStatus,
			LinkedTaskID:        task.LinkedTaskID,
			ApprovalList:        task.ApprovalList,
			Comment:             task.Comment,
			EnableDeadDateShift: task.EnableDeadDateShift,
			LayoutType:          task.LayoutType,
		})
	}

	return tasksResp, nil
}

func (u *Service) GetVacations(ctx context.Context, userID string) ([]Vacation, error) {
	vacations, err := u.db.GetTasks(ctx, userID, postgres.VacationT)
	if err != nil {
		return nil, err
	}

	vacationsResp := make([]Vacation, 0, len(vacations))
	for _, vacation := range vacations {
		vacationsResp = append(vacationsResp, Vacation{
			ObjectName:          vacation.ObjectName,
			TaskName:            vacation.TaskName,
			UUID:                vacation.UUID,
			Date:                vacation.Date,
			Author:              vacation.Author,
			AuthorID:            vacation.AuthorID,
			EndUser:             vacation.EndUser,
			EndUserID:           vacation.EndUserID,
			DeadlineDate:        vacation.DeadlineDate,
			TaskInfo:            vacation.TaskInfo,
			AgreeStatus:         vacation.AgreeStatus,
			LinkedTaskID:        vacation.LinkedTaskID,
			ApprovalList:        vacation.ApprovalList,
			Comment:             vacation.Comment,
			EnableDeadDateShift: vacation.EnableDeadDateShift,
			LayoutType:          vacation.LayoutType,
			PeriodStart:         vacation.PeriodStart,
			PeriodEnd:           vacation.PeriodEnd,
			HolidayMayker:       vacation.HolidayMayker,
			Substitutional:      vacation.Substitutional,
		})
	}

	return vacationsResp, nil
}

func (u *Service) GetPayments(ctx context.Context, userID string) ([]Payment, error) {
	payments, err := u.db.GetTasks(ctx, userID, postgres.PaymentT)
	if err != nil {
		return nil, err
	}

	paymentsResp := make([]Payment, 0, len(payments))
	for _, payment := range payments {
		paymentsResp = append(paymentsResp, Payment{
			ObjectName:          payment.ObjectName,
			TaskName:            payment.TaskName,
			UUID:                payment.UUID,
			Date:                payment.Date,
			Author:              payment.Author,
			AuthorID:            payment.AuthorID,
			EndUser:             payment.EndUser,
			EndUserID:           payment.EndUserID,
			DeadlineDate:        payment.DeadlineDate,
			TaskInfo:            payment.TaskInfo,
			AgreeStatus:         payment.AgreeStatus,
			LinkedTaskID:        payment.LinkedTaskID,
			ApprovalList:        payment.ApprovalList,
			Comment:             payment.Comment,
			EnableDeadDateShift: payment.EnableDeadDateShift,
			LayoutType:          payment.LayoutType,
			Kontragent:          payment.Kontragent,
			Organization:        payment.Organization,
			Sum:                 payment.Sum,
			PaymentDate:         payment.PaymentDate,
			PaymentPurpose:      payment.PaymentPurpose,
		})
	}

	return paymentsResp, nil
}

func (u *Service) SaveTasks(ctx context.Context, tasks []Task) error {
	tasksReq := make([]postgres.Task, 0, len(tasks))
	for _, task := range tasks {
		tasksReq = append(tasksReq, postgres.Task{
			ObjectName:          task.ObjectName,
			TaskName:            task.TaskName,
			TaskType:            postgres.TaskT,
			UUID:                task.UUID,
			Date:                task.Date,
			AuthorID:            task.AuthorID,
			CreatorID:           task.CreatorID,
			EndUserID:           task.EndUserID,
			DeadlineDate:        task.DeadlineDate,
			TaskInfo:            task.TaskInfo,
			AgreeStatus:         task.AgreeStatus,
			LinkedTaskID:        task.LinkedTaskID,
			ApprovalList:        task.ApprovalList,
			Comment:             task.Comment,
			EnableDeadDateShift: task.EnableDeadDateShift,
			LayoutType:          task.LayoutType,
		})
	}

	return u.db.CreateTasks(ctx, tasksReq)
}

func (u *Service) SaveVacations(ctx context.Context, vacations []Vacation) error {
	vacationsReq := make([]postgres.Task, 0, len(vacations))
	for _, vacation := range vacations {
		vacationsReq = append(vacationsReq, postgres.Task{
			ObjectName:          vacation.ObjectName,
			TaskName:            vacation.TaskName,
			TaskType:            postgres.VacationT,
			UUID:                vacation.UUID,
			Date:                vacation.Date,
			Author:              vacation.Author,
			AuthorID:            vacation.AuthorID,
			EndUser:             vacation.EndUser,
			EndUserID:           vacation.EndUserID,
			DeadlineDate:        vacation.DeadlineDate,
			TaskInfo:            vacation.TaskInfo,
			AgreeStatus:         vacation.AgreeStatus,
			LinkedTaskID:        vacation.LinkedTaskID,
			ApprovalList:        vacation.ApprovalList,
			Comment:             vacation.Comment,
			EnableDeadDateShift: vacation.EnableDeadDateShift,
			LayoutType:          vacation.LayoutType,
			PeriodStart:         vacation.PeriodStart,
			PeriodEnd:           vacation.PeriodEnd,
			HolidayMayker:       vacation.HolidayMayker,
			Substitutional:      vacation.Substitutional,
		})
	}

	return u.db.CreateTasks(ctx, vacationsReq)
}

func (u *Service) SavePayments(ctx context.Context, payments []Payment) error {
	paymentsReq := make([]postgres.Task, 0, len(payments))
	for _, payment := range payments {
		paymentsReq = append(paymentsReq, postgres.Task{
			ObjectName:          payment.ObjectName,
			TaskName:            payment.TaskName,
			TaskType:            postgres.PaymentT,
			UUID:                payment.UUID,
			Date:                payment.Date,
			Author:              payment.Author,
			AuthorID:            payment.AuthorID,
			EndUser:             payment.EndUser,
			EndUserID:           payment.EndUserID,
			DeadlineDate:        payment.DeadlineDate,
			TaskInfo:            payment.TaskInfo,
			AgreeStatus:         payment.AgreeStatus,
			LinkedTaskID:        payment.LinkedTaskID,
			ApprovalList:        payment.ApprovalList,
			Comment:             payment.Comment,
			EnableDeadDateShift: payment.EnableDeadDateShift,
			LayoutType:          payment.LayoutType,
			Kontragent:          payment.Kontragent,
			Organization:        payment.Organization,
			Sum:                 payment.Sum,
			PaymentDate:         payment.PaymentDate,
			PaymentPurpose:      payment.PaymentPurpose,
		})
	}

	return u.db.CreateTasks(ctx, paymentsReq)
}

func (u *Service) UpdateTaskStatus(ctx context.Context, taskID string, status bool) error {
	resp, err := u.resty.R().SetBody(map[string]any{
		"Uid":         taskID,
		"AgreeStatus": status,
	}).Post(u.basAddr + "/task/status")
	if err != nil {
		return err
	} else if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("Invalid BAS status code: %d", resp.StatusCode())
	}

	return u.db.UpdateTaskStatus(ctx, taskID, status)
}

func NewTasksService(cfg *config.Config, db DB, resty *resty.Client) *Service {
	return &Service{db: db, resty: resty, basAddr: cfg.BASAddr}
}
