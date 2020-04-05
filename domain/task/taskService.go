package task

import (
	"context"
)

type service struct {
	repository Repository
}

// NewService creates a service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Create(ctx context.Context, tsk *Task) error {
	if err := s.repository.Add(ctx, tsk); err != nil {
		return err
	}

	return nil
}

func (s *service) Get(ctx context.Context, ID string) (*Task, error) {
	tsk, err := s.repository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return tsk, nil
}

func (s *service) GetUserTask(ctx context.Context, userID string) ([]*Task, error) {
	tsks, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return tsks, nil
}
