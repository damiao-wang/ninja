package router

import(
	"fmt"
	"reflect"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
)

func AutoRouter(r *mux.Router, s interface{}) {
	vf := reflect.ValueOf(s)
	for i := 0; i < vf.NumMethod(); i++ {
		func (i int)  {
			methodName := vf.Method(i).Type().Name()
			fmt.Println("methodName: ", methodName)
			if methodName != "" && methodName != "Desc" && methodName != "Register" {
				path := fmt.Sprintf("/%v", methodName)
				r.HandleFunc(path, generateHandler(vf.Method(i))).Methods("POST")
			}
		}(i)
	}
}
	
func generateHandler(fnVal reflect.Value) http.HandlerFunc {
	fnType := fnVal.Type()
	if fnType.Kind() != reflect.Func {
		panic("fnVal is not a func.") // TODO
	}
	if fnType.NumIn() != 2 {
		panic("Num of input param isn't equal 2!") // TODO
	}
	if fnType.NumOut() != 2 {
		panic("Num of output param isn't equal 2!") // TODO
	}

	if fnType.Out(1).Name() != "error" {
		panic("The 2th output param must be error.") // TODO
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			// TODO
			panic(err)
		}
		in := reflect.New(fnType.In(0))
		err = json.Unmarshal(requestBody, in.Interface())
		if err != nil {
			panic(err) // TODO
		}
		vals := fnVal.Call([]reflect.Value{in.Elem()})
		if !vals[1].IsNil() {
			panic(vals[1].Interface()) // TODO
		}
		data, err := json.Marshal(vals[0].Interface())
		if err != nil {
			panic(err) // TODO
		}
		writer.Write(data)
	}
}