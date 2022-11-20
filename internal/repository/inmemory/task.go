package inmemory

import (
	"context"
	"oa-gogolook/internal/domain"
	"sort"
	"sync"
)

type TaskIDCounter struct {
	Mu *sync.Mutex
	ID int64
}

func NewTaskIDCounter(mu *sync.Mutex) *TaskIDCounter {
	return &TaskIDCounter{
		Mu: mu,
		ID: 0,
	}
}

func (c *TaskIDCounter) Next() int64 {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.ID++
	return c.ID
}

type TaskStore struct {
	Mu        *sync.Mutex
	IDCounter *TaskIDCounter
	tasks     map[int64]*domain.Task
}

func (t *TaskStore) AddTask(task domain.Task) (domain.Task, error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if _, ok := t.tasks[task.ID]; ok {
		return domain.Task{}, domain.ErrWrongID
	}
	t.tasks[task.ID] = &task
	return *t.tasks[task.ID], nil
}

func (t *TaskStore) DeleteTask(id int64) error {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if _, ok := t.tasks[id]; !ok {
		return domain.ErrDataNotFound
	}
	delete(t.tasks, id)
	return nil
}

func (t *TaskStore) UpdateTask(id int64, status domain.Status) (domain.Task, error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if _, ok := t.tasks[id]; !ok {
		return domain.Task{}, domain.ErrDataNotFound
	}
	task := t.tasks[id]
	task.Status = status
	return *t.tasks[id], nil
}

func NewTaskStore() *TaskStore {
	mu1 := sync.Mutex{}
	mu2 := sync.Mutex{}
	return &TaskStore{
		Mu:        &mu2,
		IDCounter: NewTaskIDCounter(&mu1),
		tasks:     map[int64]*domain.Task{},
	}
}

type taskRepository struct {
	store *TaskStore
}

func NewTaskRepository() *taskRepository {
	return &taskRepository{
		store: NewTaskStore(),
	}
}

func (r *taskRepository) List(ctx context.Context) ([]domain.Task, error) {
	var tasks []domain.Task
	for _, task := range r.store.tasks {
		tasks = append(tasks, *task)
	}
	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})
	if len(tasks) == 0 {
		return make([]domain.Task, 0), nil
	}
	return tasks, nil
}

func (r *taskRepository) Create(ctx context.Context, name string) (domain.Task, error) {
	id := r.store.IDCounter.Next()
	task := domain.Task{
		ID:     id,
		Status: domain.StatusIncomplete,
		Name:   name,
	}
	rtn, err := r.store.AddTask(task)
	if err != nil {
		return domain.Task{}, err
	}
	return rtn, nil
}

func (r *taskRepository) Update(ctx context.Context, id int64, status domain.Status) (domain.Task, error) {
	if _, ok := r.store.tasks[id]; !ok {
		return domain.Task{}, domain.ErrDataNotFound
	}
	rtn, err := r.store.UpdateTask(id, status)
	if err != nil {
		return domain.Task{}, err
	}
	return rtn, nil
}

func (r *taskRepository) Delete(ctx context.Context, id int64) error {
	err := r.store.DeleteTask(id)
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) Get(ctx context.Context, id int64) (domain.Task, error) {
	if _, ok := r.store.tasks[id]; !ok {
		return domain.Task{}, domain.ErrDataNotFound
	}
	return *r.store.tasks[id], nil
}
