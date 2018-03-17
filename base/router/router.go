package router

import (
	"fmt"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
)

func AutoRouter(r *gin.Engine, s interface{}) {
	g := r.Group(fmt.Sprintf("/api/blog"))
	fmt.Println(os.Args[0])
	vf := reflect.ValueOf(s)
	for i := 0; i < vf.NumMethod(); i++ {
		func(idx int) {
			g.POST(fmt.Sprintf("/%v", vf.Type().Method(idx).Name), func(c *gin.Context) {
				vf.Method(idx).Call([]reflect.Value{reflect.ValueOf(c)})
			})
		}(i)
	}
}
