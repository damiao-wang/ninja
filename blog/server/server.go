package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Name string
}

func (s *Service) Hello(ctx *gin.Context) {
	var name string = "wang jihang "
	ctx.String(http.StatusOK, "Hello %v.", name)
}

func (s *Service) Bye(ctx *gin.Context) {
	name, _ := ctx.Get("name")
	ctx.String(http.StatusOK, "Bye %v.", name)
}

type Handler func(*gin.Context, HelloReq) (HelloResp, error)