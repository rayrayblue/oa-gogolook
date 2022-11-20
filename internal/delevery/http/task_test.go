package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"oa-gogolook/internal/domain"
	"testing"
)

func TestTaskHandler_Create(t *testing.T) {
	expectTask := domain.CreateTaskResponse{Result: domain.Task{
		ID:     1,
		Status: 0,
		Name:   "TaskName1",
	}}
	tests := []struct {
		name          string
		query         domain.CreateTaskRequest
		buildStubs    func(s *domain.TaskUseCase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: domain.CreateTaskRequest{
				Name: expectTask.Result.Name,
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)
				var task domain.CreateTaskResponse
				err = json.Unmarshal(data, &task)
				require.NoError(t, err)
				require.Equal(t, expectTask, task)
			},
		},
		{
			name: "EmptyName",
			query: domain.CreateTaskRequest{
				Name: "",
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newTestServer(t)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/task")
			data, err := json.Marshal(tt.query)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.Router.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestTaskHandler_List(t *testing.T) {
	tests := []struct {
		name          string
		buildStubs    func(s *domain.TaskUseCase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(s *domain.TaskUseCase) {
				var u domain.TaskUseCase
				u = *s
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: "TaskName1"})
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: "TaskName2"})
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: "TaskName3"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)
				var tasks domain.ListTaskResponse
				err = json.Unmarshal(data, &tasks)
				require.NoError(t, err)
				require.Equal(t, 3, len(tasks.Result))
			},
		},
		{
			name: "EmptyTasks",
			buildStubs: func(s *domain.TaskUseCase) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)
				var tasks domain.ListTaskResponse
				err = json.Unmarshal(data, &tasks)
				require.NoError(t, err)
				require.Equal(t, 0, len(tasks.Result))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newTestServer(t)
			tt.buildStubs(&server.U)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/tasks")
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.Router.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestTaskHandler_Delete(t *testing.T) {
	tests := []struct {
		name          string
		deleteID      int64
		buildStubs    func(s *domain.TaskUseCase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, s *domain.TaskUseCase)
	}{
		{
			name: "OK",
			buildStubs: func(s *domain.TaskUseCase) {
				var u domain.TaskUseCase
				u = *s
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: "TaskName1"})
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: "TaskName2"})
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: "TaskName3"})
			},
			deleteID: 3,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, s *domain.TaskUseCase) {
				var u domain.TaskUseCase
				u = *s
				require.Equal(t, http.StatusOK, recorder.Code)
				_, err := u.Get(context.Background(), 3)
				require.ErrorIs(t, err, domain.ErrDataNotFound)
			},
		},
		{
			name: "NotFound",
			buildStubs: func(s *domain.TaskUseCase) {
			},
			deleteID: 3,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, s *domain.TaskUseCase) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "BadParam",
			buildStubs: func(s *domain.TaskUseCase) {
			},
			deleteID: 0,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, s *domain.TaskUseCase) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newTestServer(t)
			tt.buildStubs(&server.U)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/task/%d", tt.deleteID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)
			server.Router.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder, &server.U)
		})
	}
}

func TestTaskHandler_Update(t *testing.T) {
	taskID := int64(1)
	taskName := "TaskName1"
	tests := []struct {
		name          string
		taskID        int64
		query         domain.UpdateTaskRequest
		buildStubs    func(s *domain.TaskUseCase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, s *domain.TaskUseCase)
	}{
		{
			name:   "OK",
			taskID: taskID,
			query: domain.UpdateTaskRequest{
				Name:   taskName,
				ID:     taskID,
				Status: &domain.StatusComplete,
			},
			buildStubs: func(s *domain.TaskUseCase) {
				var u domain.TaskUseCase
				u = *s
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: taskName})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, s *domain.TaskUseCase) {
				var u domain.TaskUseCase
				u = *s
				require.Equal(t, http.StatusOK, recorder.Code)
				task, err := u.Get(context.Background(), taskID)
				require.NoError(t, err)
				require.Equal(t, taskName, task.Name)
				require.Equal(t, domain.StatusComplete, task.Status)
			},
		},
		{
			name:   "InvalidParamID",
			taskID: 0,
			query: domain.UpdateTaskRequest{
				Name:   taskName,
				ID:     taskID,
				Status: &domain.StatusComplete,
			},
			buildStubs: func(s *domain.TaskUseCase) {
				var u domain.TaskUseCase
				u = *s
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: taskName})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, s *domain.TaskUseCase) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InvalidRequestPayloda",
			taskID: taskID,
			query: domain.UpdateTaskRequest{
				Name:   taskName,
				ID:     0,
				Status: &domain.StatusComplete,
			},
			buildStubs: func(s *domain.TaskUseCase) {
				var u domain.TaskUseCase
				u = *s
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: taskName})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, s *domain.TaskUseCase) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "StatusNoChange",
			taskID: taskID,
			query: domain.UpdateTaskRequest{
				Name:   taskName,
				ID:     taskID,
				Status: &domain.StatusIncomplete,
			},
			buildStubs: func(s *domain.TaskUseCase) {
				var u domain.TaskUseCase
				u = *s
				_, _ = u.Create(context.Background(), domain.CreateTaskRequest{Name: taskName})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, s *domain.TaskUseCase) {
				var u domain.TaskUseCase
				u = *s
				require.Equal(t, http.StatusOK, recorder.Code)
				task, err := u.Get(context.Background(), taskID)
				require.NoError(t, err)
				require.Equal(t, taskName, task.Name)
				require.Equal(t, domain.StatusIncomplete, task.Status)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newTestServer(t)
			tt.buildStubs(&server.U)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/task/%d", tt.taskID)
			data, err := json.Marshal(tt.query)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.Router.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder, &server.U)
		})
	}
}
