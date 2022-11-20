package main

import (
	"log"
	"oa-gogolook/internal"
	"oa-gogolook/internal/domain"
	"oa-gogolook/internal/repository/inmemory"
	"oa-gogolook/internal/usecase"
)

func main() {
	r := inmemory.NewTaskRepository()
	u := usecase.NewTaskUsecase(r)
	config, err := domain.LoadConfig("configs", "app.env")
	if err != nil {
		log.Fatal("can not load config. ", err)
	}

	server, err := internal.NewHttpServer(u, config)
	if err != nil {
		log.Fatal("can not create server ", err.Error())
	}
	server.Start()
}
