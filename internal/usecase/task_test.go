package usecase

import (
	"context"
	"oa-gogolook/internal/domain"
	"oa-gogolook/internal/repository/inmemory"
	"reflect"
	"testing"
)

func Test_taskUsecase_Create(t *testing.T) {
	type fields struct {
		taskRepository domain.TaskRepository
	}
	type args struct {
		ctx context.Context
		req domain.CreateTaskRequest
	}
	tests := []struct {
		name       string
		buildStubs func(repo *domain.TaskRepository)
		fields     fields
		args       args
		want       domain.CreateTaskResponse
		wantErr    bool
	}{
		{
			name: "OK",
			buildStubs: func(repo *domain.TaskRepository) {

			},
			fields: fields{
				taskRepository: inmemory.NewTaskRepository(),
			},
			args: args{
				ctx: context.Background(),
				req: domain.CreateTaskRequest{
					Name: "taskName",
				},
			},
			want: domain.CreateTaskResponse{
				Result: domain.Task{
					ID:     1,
					Status: domain.StatusIncomplete,
					Name:   "taskName",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &taskUsecase{
				taskRepository: tt.fields.taskRepository,
			}
			got, err := u.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskUsecase_List(t *testing.T) {
	type fields struct {
		taskRepository domain.TaskRepository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		buildStubs func(repo domain.TaskRepository)
		fields     fields
		args       args
		want       domain.ListTaskResponse
		wantErr    bool
	}{
		{
			name: "OK",
			buildStubs: func(repo domain.TaskRepository) {
				_, _ = repo.Create(context.Background(), "taskName1")
				_, _ = repo.Create(context.Background(), "taskName2")
				_, _ = repo.Create(context.Background(), "taskName3")
			},
			fields: fields{
				taskRepository: inmemory.NewTaskRepository(),
			},
			args: args{
				ctx: context.Background(),
			},
			want: domain.ListTaskResponse{
				Result: []domain.Task{
					{
						ID:     1,
						Status: domain.StatusIncomplete,
						Name:   "taskName1",
					},
					{
						ID:     2,
						Status: domain.StatusIncomplete,
						Name:   "taskName2",
					},
					{
						ID:     3,
						Status: domain.StatusIncomplete,
						Name:   "taskName3",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "OKEmpty",
			buildStubs: func(repo domain.TaskRepository) {
			},
			fields: fields{
				taskRepository: inmemory.NewTaskRepository(),
			},
			args: args{
				ctx: context.Background(),
			},
			want: domain.ListTaskResponse{
				Result: []domain.Task{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &taskUsecase{
				taskRepository: tt.fields.taskRepository,
			}
			tt.buildStubs(u.taskRepository)
			got, err := u.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				if len(got.Result) == 0 && len(tt.want.Result) == 0 {
					return
				}
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskUsecase_Update(t *testing.T) {
	type fields struct {
		taskRepository domain.TaskRepository
	}
	type args struct {
		ctx context.Context
		req domain.UpdateTaskRequest
	}
	tests := []struct {
		name       string
		buildStubs func(repo domain.TaskRepository)
		fields     fields
		args       args
		want       domain.UpdateTaskResponse
		wantErr    bool
	}{
		{
			name: "OK",
			buildStubs: func(repo domain.TaskRepository) {
				_, _ = repo.Create(context.Background(), "taskName1")
				_, _ = repo.Create(context.Background(), "taskName2")
				_, _ = repo.Create(context.Background(), "taskName3")
			},
			fields: fields{taskRepository: inmemory.NewTaskRepository()},
			args: args{
				ctx: context.Background(),
				req: domain.UpdateTaskRequest{
					ID:     1,
					Status: &domain.StatusComplete,
					Name:   "taskName1",
				},
			},
			want: domain.UpdateTaskResponse{
				Result: domain.Task{
					ID:     1,
					Status: domain.StatusComplete,
					Name:   "taskName1",
				},
			},
			wantErr: false,
		},
		{
			name: "NameNotMatch",
			buildStubs: func(repo domain.TaskRepository) {
				_, _ = repo.Create(context.Background(), "taskName1")
				_, _ = repo.Create(context.Background(), "taskName2")
				_, _ = repo.Create(context.Background(), "taskName3")
			},
			fields: fields{taskRepository: inmemory.NewTaskRepository()},
			args: args{
				ctx: context.Background(),
				req: domain.UpdateTaskRequest{
					ID:     1,
					Status: &domain.StatusComplete,
					Name:   "taskNameX",
				},
			},
			want:    domain.UpdateTaskResponse{},
			wantErr: true,
		},
		{
			name: "TaskNotFound",
			buildStubs: func(repo domain.TaskRepository) {
				_, _ = repo.Create(context.Background(), "taskName1")
				_, _ = repo.Create(context.Background(), "taskName2")
				_, _ = repo.Create(context.Background(), "taskName3")
			},
			fields: fields{taskRepository: inmemory.NewTaskRepository()},
			args: args{
				ctx: context.Background(),
				req: domain.UpdateTaskRequest{
					ID:     5,
					Status: &domain.StatusComplete,
					Name:   "taskName1",
				},
			},
			want:    domain.UpdateTaskResponse{},
			wantErr: true,
		},
		{
			name: "StatusNotChange",
			buildStubs: func(repo domain.TaskRepository) {
				_, _ = repo.Create(context.Background(), "taskName1")
				_, _ = repo.Create(context.Background(), "taskName2")
				_, _ = repo.Create(context.Background(), "taskName3")
			},
			fields: fields{taskRepository: inmemory.NewTaskRepository()},
			args: args{
				ctx: context.Background(),
				req: domain.UpdateTaskRequest{
					ID:     1,
					Status: &domain.StatusIncomplete,
					Name:   "taskName1",
				},
			},
			want: domain.UpdateTaskResponse{
				Result: domain.Task{
					ID:     1,
					Status: domain.StatusIncomplete,
					Name:   "taskName1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &taskUsecase{
				taskRepository: tt.fields.taskRepository,
			}
			tt.buildStubs(u.taskRepository)
			got, err := u.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskUsecase_Delete(t *testing.T) {
	type fields struct {
		taskRepository domain.TaskRepository
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name       string
		buildStubs func(repo domain.TaskRepository)
		fields     fields
		args       args
		wantErr    bool
	}{
		{
			name: "OK",
			buildStubs: func(repo domain.TaskRepository) {
				_, _ = repo.Create(context.Background(), "taskName1")
				_, _ = repo.Create(context.Background(), "taskName2")
				_, _ = repo.Create(context.Background(), "taskName3")
			},
			fields: fields{taskRepository: inmemory.NewTaskRepository()},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: false,
		},
		{
			name: "NotFound",
			buildStubs: func(repo domain.TaskRepository) {
				_, _ = repo.Create(context.Background(), "taskName1")
				_, _ = repo.Create(context.Background(), "taskName2")
				_, _ = repo.Create(context.Background(), "taskName3")
			},
			fields: fields{taskRepository: inmemory.NewTaskRepository()},
			args: args{
				ctx: context.Background(),
				id:  5,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &taskUsecase{
				taskRepository: tt.fields.taskRepository,
			}
			tt.buildStubs(u.taskRepository)
			if err := u.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskUsecase_Get(t *testing.T) {
	type fields struct {
		taskRepository domain.TaskRepository
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name       string
		buildStubs func(repo domain.TaskRepository)
		fields     fields
		args       args
		want       domain.Task
		wantErr    bool
	}{
		{
			name: "OK",
			buildStubs: func(repo domain.TaskRepository) {
				_, _ = repo.Create(context.Background(), "taskName1")
				_, _ = repo.Create(context.Background(), "taskName2")
				_, _ = repo.Create(context.Background(), "taskName3")
			},
			fields: fields{taskRepository: inmemory.NewTaskRepository()},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: domain.Task{
				ID:     1,
				Status: domain.StatusIncomplete,
				Name:   "taskName1",
			},
			wantErr: false,
		},
		{
			name: "NotFound",
			buildStubs: func(repo domain.TaskRepository) {
				_, _ = repo.Create(context.Background(), "taskName1")
				_, _ = repo.Create(context.Background(), "taskName2")
				_, _ = repo.Create(context.Background(), "taskName3")
			},
			fields: fields{taskRepository: inmemory.NewTaskRepository()},
			args: args{
				ctx: context.Background(),
				id:  5,
			},
			want:    domain.Task{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &taskUsecase{
				taskRepository: tt.fields.taskRepository,
			}
			tt.buildStubs(u.taskRepository)
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
