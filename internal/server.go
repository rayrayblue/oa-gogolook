package internal

import (
	"github.com/gin-gonic/gin"
	"log"
	"oa-gogolook/internal/delevery/http"
	"oa-gogolook/internal/domain"
)

type Server struct {
	Router *gin.Engine
	config domain.AppConfig
}

func NewHttpServer(usecase domain.TaskUseCase, config domain.AppConfig) (*Server, error) {
	router := gin.Default()
	http.NewTaskHandler(router, usecase)
	server := &Server{}
	server.Router = router
	server.config = config
	return server, nil
}

func (server *Server) Start() {
	if server.config.ServerAddress == "" {
		log.Fatal("server address is empty")
	}
	err := server.Router.Run(server.config.ServerAddress)
	if err != nil {
		panic(err)
	}
}
