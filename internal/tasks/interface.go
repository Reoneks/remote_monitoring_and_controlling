package tasks

import (
	"context"
	"remote_monitoring_and_controlling/pkg/postgres"
)

type DB interface {
	GetTasks(ctx context.Context, userID string, taskType postgres.TaskType) ([]postgres.Task, error)
	CreateTasks(ctx context.Context, tasks []postgres.Task) error
	UpdateTaskStatus(ctx context.Context, taskID string, status bool) error
}
