package usecase

import (
	"context"
	"oa-gogolook/internal/domain"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
}

func NewTaskUsecase(taskRepository domain.TaskRepository) *taskUsecase {
	return &taskUsecase{
		taskRepository: taskRepository,
	}
}

func (u *taskUsecase) List(ctx context.Context) (domain.ListTaskResponse, error) {
	var rtn domain.ListTaskResponse
	got, err := u.taskRepository.List(ctx)
	if err != nil {
		return rtn, err
	}
	rtn.Result = got
	return rtn, nil
}

func (u *taskUsecase) Create(ctx context.Context, req domain.CreateTaskRequest) (domain.CreateTaskResponse, error) {
	var rtn domain.CreateTaskResponse
	got, err := u.taskRepository.Create(ctx, req.Name)
	if err != nil {
		return rtn, err
	}
	rtn.Result = got
	return rtn, nil

}

func (u *taskUsecase) Update(ctx context.Context, req domain.UpdateTaskRequest) (domain.UpdateTaskResponse, error) {
	var rtn domain.UpdateTaskResponse
	gotTask, err := u.taskRepository.Get(ctx, req.ID)
	if err != nil {
		return rtn, err
	}
	if gotTask.Name != req.Name {
		return rtn, domain.ErrTaskNameNotMatch
	}

	got, err := u.taskRepository.Update(ctx, req.ID, *req.Status)
	if err != nil {
		return rtn, err
	}
	rtn.Result = got
	return rtn, nil
}

func (u *taskUsecase) Delete(ctx context.Context, id int64) error {
	err := u.taskRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *taskUsecase) Get(ctx context.Context, id int64) (domain.Task, error) {
	got, err := u.taskRepository.Get(ctx, id)
	if err != nil {
		return got, err
	}
	return got, nil
}
