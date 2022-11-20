package inmemory

import (
	"context"
	"oa-gogolook/internal/domain"
	"reflect"
	"testing"
)

func Test_taskRepository_Create(t *testing.T) {
	type fields struct {
		store *TaskStore
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name       string
		buildStubs func(store *TaskStore)
		fields     fields
		args       args
		want       domain.Task
		wantErr    bool
	}{
		{
			name: "OK",
			buildStubs: func(store *TaskStore) {

			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx:  context.Background(),
				name: "taskName",
			},
			want: domain.Task{
				ID:     1,
				Status: domain.StatusIncomplete,
				Name:   "taskName",
			},
			wantErr: false,
		},
		{
			name: "WrongID",
			buildStubs: func(store *TaskStore) {
				store.tasks[1] = &domain.Task{
					ID:     1,
					Status: domain.StatusIncomplete,
					Name:   "taskName1",
				}
			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx:  context.Background(),
				name: "taskName",
			},
			want:    domain.Task{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &taskRepository{
				store: tt.fields.store,
			}

			tt.buildStubs(r.store)

			got, err := r.Create(tt.args.ctx, tt.args.name)

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

func Test_taskRepository_List(t *testing.T) {
	type fields struct {
		store *TaskStore
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		buildStubs func(store *TaskStore)
		fields     fields
		args       args
		want       []domain.Task
		wantErr    bool
	}{
		{
			name: "OK",
			buildStubs: func(store *TaskStore) {
				id := store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusIncomplete,
					Name:   "taskName1",
				}
				id = store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusIncomplete,
					Name:   "taskName2",
				}
			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx: context.Background(),
			},
			want: []domain.Task{
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
			},
			wantErr: false,
		},
		{
			name: "OKEmpty",
			buildStubs: func(store *TaskStore) {
			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []domain.Task{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &taskRepository{
				store: tt.fields.store,
			}
			tt.buildStubs(r.store)
			got, err := r.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				if len(got) == 0 && len(tt.want) == 0 {
					return
				}
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskRepository_Update(t *testing.T) {
	type fields struct {
		store *TaskStore
	}
	type args struct {
		ctx    context.Context
		id     int64
		status domain.Status
	}
	tests := []struct {
		name       string
		buildStubs func(store *TaskStore)
		fields     fields
		args       args
		want       domain.Task
		wantErr    bool
	}{
		{
			name: "OK",
			buildStubs: func(store *TaskStore) {
				id := store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusIncomplete,
					Name:   "taskName1",
				}
			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx:    context.Background(),
				id:     1,
				status: domain.StatusComplete,
			},
			want: domain.Task{
				ID:     1,
				Status: domain.StatusComplete,
				Name:   "taskName1",
			},
			wantErr: false,
		},
		{
			name: "NoChange",
			buildStubs: func(store *TaskStore) {
				id := store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusIncomplete,
					Name:   "taskName1",
				}
			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx:    context.Background(),
				id:     1,
				status: domain.StatusIncomplete,
			},
			want: domain.Task{
				ID:     1,
				Status: domain.StatusIncomplete,
				Name:   "taskName1",
			},
			wantErr: false,
		},
		{
			name: "NotExists",
			buildStubs: func(store *TaskStore) {
			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx:    context.Background(),
				id:     1,
				status: domain.StatusIncomplete,
			},
			want:    domain.Task{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &taskRepository{
				store: tt.fields.store,
			}
			tt.buildStubs(r.store)
			got, err := r.Update(tt.args.ctx, tt.args.id, tt.args.status)
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

func Test_taskRepository_Delete(t *testing.T) {
	type fields struct {
		store *TaskStore
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name       string
		buildStubs func(store *TaskStore)
		fields     fields
		args       args
		wantErr    bool
	}{
		{
			name: "OK",
			buildStubs: func(store *TaskStore) {
				id := store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusIncomplete,
					Name:   "taskName1",
				}
			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: false,
		},
		{
			name: "NotExists",
			buildStubs: func(store *TaskStore) {
			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &taskRepository{
				store: tt.fields.store,
			}
			tt.buildStubs(r.store)
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskRepository_Get(t *testing.T) {
	type fields struct {
		store *TaskStore
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name       string
		buildStubs func(store *TaskStore)
		fields     fields
		args       args
		want       domain.Task
		wantErr    bool
	}{
		{
			name: "OKIncompleteTask",
			buildStubs: func(store *TaskStore) {
				id := store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusIncomplete,
					Name:   "taskName1",
				}
				id = store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusComplete,
					Name:   "taskName2",
				}
				id = store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusComplete,
					Name:   "taskName3",
				}
			},
			fields: fields{
				store: NewTaskStore(),
			},
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
			name: "OKCompleteTask",
			buildStubs: func(store *TaskStore) {
				id := store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusIncomplete,
					Name:   "taskName1",
				}
				id = store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusComplete,
					Name:   "taskName2",
				}
				id = store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusComplete,
					Name:   "taskName3",
				}
			},
			fields: fields{
				store: NewTaskStore(),
			},
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			want: domain.Task{
				ID:     2,
				Status: domain.StatusComplete,
				Name:   "taskName2",
			},
			wantErr: false,
		},
		{
			name: "NotExists",
			buildStubs: func(store *TaskStore) {
				id := store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusIncomplete,
					Name:   "taskName1",
				}
				id = store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusComplete,
					Name:   "taskName2",
				}
				id = store.IDCounter.Next()
				store.tasks[id] = &domain.Task{
					ID:     id,
					Status: domain.StatusComplete,
					Name:   "taskName3",
				}
			},
			fields: fields{
				store: NewTaskStore(),
			},
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
			r := &taskRepository{
				store: tt.fields.store,
			}
			tt.buildStubs(r.store)
			got, err := r.Get(tt.args.ctx, tt.args.id)
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
