package task

import (
	"errors"
	"time"
)

// Lit bit of ddd modeling

// Status defines task valid status
type Status string

// Status all type
const (
	Pending  Status = "Pending"
	Active   Status = "Active"
	InActive Status = "InActive"
	Complete Status = "Complete"
)

// Task defines Task type
type Task struct {
	ID          string
	UserID      string
	Topic       string
	Description string
	Status      Status
	SubTasks    []*SubTask
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

// Complete makes the task status Complete
func (t *Task) Complete() error {
	if t.Status == Complete || t.Status == InActive {
		return errors.New("Can't complete a inActvie Or complete task")
	}

	t.Status = Complete
	return nil
}

// AddSubTask adds subTask to a task
func (t *Task) AddSubTask(sbTsk *SubTask) error {
	if t.Status == InActive || t.Status == Complete {
		return errors.New("Can't add subTask to a inActvie Or complete task")
	}

	t.SubTasks = append(t.SubTasks, sbTsk)
	return nil
}
