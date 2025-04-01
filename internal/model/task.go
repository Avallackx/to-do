package model

import (
	"context"
	"time"
	"todo-app/internal/utils"
)

type TaskRepository interface {
	Create(ctx context.Context, input *Task) (err error)
	DeleteByID(ctx context.Context, ID int64) (err error)
	FindByID(ctx context.Context, ID int64) (task *Task, err error)
	FindAll(ctx context.Context) (tasks []*Task, err error)
	Update(ctx context.Context, input *Task) (task *Task, err error)
}

type TaskUsecase interface {
	Create(ctx context.Context, input *Task) (task *Task, err error)
	DeleteByID(ctx context.Context, ID int64) (err error)
	FindByID(ctx context.Context, ID int64) (task *Task, err error)
	FindAll(ctx context.Context) (tasks []*Task, err error)
	Update(ctx context.Context, input *Task) (task *Task, err error)
}

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Todo      string    `json:"todo"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateTaskInput struct {
	Title     string `json:"title"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
}

func (i CreateTaskInput) ToModel() *Task {
	return &Task{
		ID:        utils.GenerateID(),
		Title:     i.Title,
		Todo:      i.Todo,
		Completed: i.Completed,
		CreatedAt: time.Now(),
	}
}

type UpdateTaskInput struct {
	Title     string `json:"title"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
}

type GetTasksQueryParams struct {
	Page int64 `query:"page"`
	Size int64 `query:"size"`
}

func (i UpdateTaskInput) ToModel() *Task {
	return &Task{
		Title:     i.Title,
		Todo:      i.Todo,
		Completed: i.Completed,
		UpdatedAt: time.Now(),
	}
}
