package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/BackAged/go-hexagonal-architecture/domain/task"
	"github.com/BackAged/go-hexagonal-architecture/infrastructure/database"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
)

type jsonTask struct {
	ID          string          `json:"ID"`
	UserID      string          `json:"userID"`
	Topic       string          `json:"topic"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	SubTasks    []*task.SubTask `json:"sub_task"`
	CreatedAt   *time.Time      `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at"`
}

func toJSON(tsk *task.Task) *jsonTask {
	return &jsonTask{
		ID:          tsk.ID,
		UserID:      tsk.UserID,
		Topic:       tsk.Topic,
		Description: tsk.Description,
		SubTasks:    tsk.SubTasks,
		Status:      string(tsk.Status),
		CreatedAt:   tsk.CreatedAt,
		UpdatedAt:   tsk.UpdatedAt,
	}
}

func toModel(tsk *jsonTask) *task.Task {
	return &task.Task{
		ID:          tsk.ID,
		UserID:      tsk.UserID,
		Topic:       tsk.Topic,
		Description: tsk.Description,
		SubTasks:    tsk.SubTasks,
		Status:      task.Status(tsk.Status),
		CreatedAt:   tsk.CreatedAt,
		UpdatedAt:   tsk.UpdatedAt,
	}
}

type taskRepository struct {
	client *database.InMemoryClient
}

// NewTaskRepository returns a new taskRepository
func NewTaskRepository(client *database.InMemoryClient) task.Repository {
	return &taskRepository{
		client: client,
	}
}

// Add adds into repository
func (tr *taskRepository) makeKey(tskID string, usrID string) string {
	return fmt.Sprintf("task.%s.%s", tskID, usrID)
}

// Add adds into repository
func (tr *taskRepository) Add(ctx context.Context, tsk *task.Task) error {
	if tsk.ID == "" {
		tsk.ID = uuid.New().String()
	}
	key := tr.makeKey(tsk.ID, tsk.UserID)

	now := time.Now()
	tsk.CreatedAt = &now
	tsk.UpdatedAt = &now

	jsnTsk, err := json.Marshal(toJSON(tsk))
	if err != nil {
		return err
	}

	err = tr.client.Insert(ctx, key, jsnTsk, 0)
	if err != nil {
		return err
	}

	return nil
}

// FindByID finds something by ID
func (tr *taskRepository) FindByID(ctx context.Context, tskID string) (*task.Task, error) {
	key := tr.makeKey(tskID, "*")

	tsk := &jsonTask{}

	rows := tr.client.Find(ctx, key)
	if rows.Next() {
		if err := rows.Scan(&tsk); err != nil {
			fmt.Println("here", err)
			return nil, err
		}
	} else {
		return nil, nil
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	return toModel(tsk), nil
}

// FindByUserID finds by userID
func (tr *taskRepository) FindByUserID(ctx context.Context, usrID string) ([]*task.Task, error) {
	key := tr.makeKey("*", usrID)

	tsks := []*task.Task{}

	rows := tr.client.Find(ctx, key)
	for rows.Next() {
		tsk := &jsonTask{}
		if err := rows.Scan(&tsk); err != nil {
			return nil, err
		}
		tsks = append(tsks, toModel(tsk))
	}

	if err := rows.Err(); err != nil {
		return tsks, err
	}

	return tsks, nil
}
