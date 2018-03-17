package main

import (
	"ninja/base/router"
	"ninja/blog/server"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.AutoRouter(r, &server.Service{})
	// r.Run(":8080")
	//r.POST("/api/blog/Hello", s.Hello)
	r.Run(":8080")
}
