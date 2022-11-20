package http

import (
	"github.com/gin-gonic/gin"
	"oa-gogolook/internal/domain"
	"oa-gogolook/internal/repository/inmemory"
	"oa-gogolook/internal/usecase"
	"os"
	"testing"
)

type TestServer struct {
	Router *gin.Engine
	U      domain.TaskUseCase
}

func newTestServer(t *testing.T) TestServer {
	r := inmemory.NewTaskRepository()
	u := usecase.NewTaskUsecase(r)

	router := gin.Default()
	NewTaskHandler(router, u)
	server := TestServer{
		Router: router,
		U:      u,
	}
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
