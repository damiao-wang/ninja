package router

import (
	"fmt"
	"reflect"

	"ninja/base/misc/stack"

	"github.com/gin-gonic/gin"
)

func AutoRouter(r *gin.Engine, s interface{}) {
	reflectVal := reflect.ValueOf(s)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	prefix := fmt.Sprintf("/api/%v", stack.PkgName(0))
	fmt.Println("prefix: ", prefix)
	for i := 0; i < rt.NumMethod(); i++ {
		r.POST(fmt.Sprintf("%v/%v", prefix, ct.Method(i).Name), rt.Method(i).Func.Interface().(gin.HandlerFunc))
	}
}
