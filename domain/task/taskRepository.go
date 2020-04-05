package task

import "context"

// Repository defines the port of task for infrastracture adapter
type Repository interface {
	Add(context.Context, *Task) error
	FindByID(context.Context, string) (*Task, error)
	FindByUserID(context.Context, string) ([]*Task, error)
}
