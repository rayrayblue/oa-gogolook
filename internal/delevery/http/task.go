package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oa-gogolook/internal/domain"
)

type TaskHandler struct {
	taskUsecse domain.TaskUseCase
}

func NewTaskHandler(e *gin.Engine, taskUsecse domain.TaskUseCase) {
	h := &TaskHandler{
		taskUsecse: taskUsecse,
	}
	e.GET("/tasks", h.List)
	e.POST("/task", h.Create)
	e.PUT("/task/:task_id", h.Update)
	e.DELETE("/task/:task_id", h.Delete)
}

func (h *TaskHandler) Create(ctx *gin.Context) {
	var req domain.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrInvalidPayload)
		return
	}

	rtn, err := h.taskUsecse.Create(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, rtn)
}

func (h *TaskHandler) List(ctx *gin.Context) {
	rtn, err := h.taskUsecse.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, rtn)
}

func (h *TaskHandler) Delete(ctx *gin.Context) {
	var req domain.DeleteTaskRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrInvalidParameters)
		return
	}
	err := h.taskUsecse.Delete(ctx, req.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (h *TaskHandler) Update(ctx *gin.Context) {
	var para domain.UpdateTaskUriParameter
	var req domain.UpdateTaskRequest
	if err := ctx.ShouldBindUri(&para); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrInvalidParameters)
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrInvalidPayload)
		return
	}
	if para.ID != req.ID {
		ctx.JSON(http.StatusBadRequest, domain.ErrInvalidParameters)
		return
	}

	rtn, err := h.taskUsecse.Update(ctx, req)
	if err != nil {
		if err == domain.ErrDataNotFound {
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, rtn)
}
