package task

import "context"

// Service provides port for application adapter.
type Service interface {
	Create(context.Context, *Task) error
	Get(context.Context, string) (*Task, error)
	GetUserTask(context.Context, string) ([]*Task, error)
}
