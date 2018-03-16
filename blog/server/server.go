package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service struct{}

func (s *Service) Hello(ctx *gin.Context) {
	name := ctx.PostForm("name")
	ctx.String(http.StatusOK, "Hello %v.", name)
}

func (s *Service) Bye(ctx *gin.Context) {
	name := ctx.PostForm("name")
	ctx.String(http.StatusOK, "Bye %v.", name)
}
