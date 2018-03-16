package main

import (
	"ninja/base/router"
	"ninja/blog/server"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	s := &server.Service{}
	r.POST("/api/blog/Hello", s.Hello)
	router.AutoRouter(r, &server.Service{})
	r.Run("8080")
}
