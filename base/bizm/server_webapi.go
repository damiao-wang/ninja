package bizm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"
)

type WebServer struct {
	mux        *mux.Router
	middleware []negroni.Handler
}

func (s *WebServer) AddMiddleware(m ...negroni.Handler) {
	s.middleware = append(s.middleware, m...)
}

func (s *WebServer) Serve(ln net.Listener) error {
	if s.mux == nil {
		s.mux = mux.NewRouter()
	}

	n := negroni.Classic().With(s.middleware...)
	n.UseHandler(s.mux)
	srv := http.Server{
		Handler: n,
	}
	return srv.Serve(ln)
}

func (s *WebServer) AutoRouter(c interface{}) {
	if s.mux == nil {
		s.mux = mux.NewRouter()
	}

	serviceName := getServiceName(c)
	vf := reflect.ValueOf(c)
	ctx := context.Background()
	for i := 0; i < vf.NumMethod(); i++ {
		func(i int) {
			path := fmt.Sprintf("/api/%v/%v", serviceName, vf.Type().Method(i).Name)
			s.mux.HandleFunc(path, generateHandler(ctx, vf.Method(i))).Methods("POST")
		}(i)
	}
}

func generateHandler(ctx context.Context, fnVal reflect.Value) http.HandlerFunc {
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
		in := reflect.New(fnType.In(1))
		err = json.Unmarshal(requestBody, in.Interface())
		if err != nil {
			panic(err) // TODO
		}
		vals := fnVal.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			in.Elem(),
		})
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

func getServiceName(s interface{}) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return path.Base(t.PkgPath())
}
