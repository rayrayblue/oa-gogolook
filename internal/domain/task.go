package domain

import (
	"context"
)

type Status int64

var StatusIncomplete Status = 0
var StatusComplete Status = 1

type Task struct {
	ID     int64  `json:"id"`
	Status Status `json:"status"`
	Name   string `json:"name"`
}

type CreateTaskRequest struct {
	Name string `json:"name" binding:"required" `
}

type CreateTaskResponse struct {
	Result Task `json:"result"`
}

type UpdateTaskUriParameter struct {
	ID int64 `uri:"task_id" binding:"required,min=1"`
}

type UpdateTaskRequest struct {
	ID     int64   `json:"id" binding:"required"`
	Status *Status `json:"status" binding:"required,min=0,max=1"`
	Name   string  `json:"name" binding:"required"`
}

type UpdateTaskResponse struct {
	Result Task `json:"result"`
}

type ListTaskResponse struct {
	Result []Task `json:"result"`
}

type DeleteTaskRequest struct {
	ID int64 `uri:"task_id" binding:"required,min=1"`
}

type TaskUseCase interface {
	List(ctx context.Context) (ListTaskResponse, error)
	Create(ctx context.Context, req CreateTaskRequest) (CreateTaskResponse, error)
	Update(ctx context.Context, req UpdateTaskRequest) (UpdateTaskResponse, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (Task, error)
}

type TaskRepository interface {
	List(ctx context.Context) ([]Task, error)
	Create(ctx context.Context, name string) (Task, error)
	Update(ctx context.Context, id int64, status Status) (Task, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (Task, error)
}
