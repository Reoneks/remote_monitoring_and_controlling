package postgres

import "context"

func (p *Postgres) GetTasks(ctx context.Context, userID string, taskType TaskType) ([]Task, error) {
	var result []Task
	return result, p.db.WithContext(ctx).Model(&Task{}).Select("*").Where("end_user_id = ? AND task_type = ?", userID, taskType).Find(&result).Error
}

func (p *Postgres) CreateTasks(ctx context.Context, tasks []Task) error {
	return p.db.Model(&Task{}).WithContext(ctx).Create(tasks).Error
}

func (p *Postgres) UpdateTaskStatus(ctx context.Context, taskID string, status bool) error {
	return p.db.Model(&Task{}).WithContext(ctx).Where("uuid = ?", taskID).Update("agree_status", status).Error
}
