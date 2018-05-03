package bizm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"reflect"
	"strings"
	"time"

	"ninja/base/mconf"
	"ninja/base/misc/context"
	"ninja/base/misc/errors"
	"ninja/base/misc/grace"
	"ninja/base/misc/log"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type WebServer struct {
	srv        *http.Server
	mux        *mux.Router
	middleware []negroni.Handler
	IsDownload bool
}

type Filer interface {
	GetFilename() string
	GetData() []byte
	GetContentType() string
}

func (s *WebServer) AddMiddleware(m ...negroni.Handler) {
	s.middleware = append(s.middleware, m...)
}

func (s *WebServer) Serve(ln net.Listener) {
	if s.mux == nil {
		s.mux = mux.NewRouter()
	}

	n := negroni.Classic()
	n.With(s.middleware...)
	n.UseHandler(s.mux)
	s.srv = &http.Server{
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	grace.CountInTask()
	log.Infof("server HTTP start. Listen: %v", ln.Addr())
	s.srv.Serve(ln)
	return
}

func (s *WebServer) RegisterRouter(mux *mux.Router) {
	if mux != nil {
		s.mux = mux
	}
}

func (s *WebServer) GenHTTPHandler(fn interface{}) http.HandlerFunc {
	fnVal := reflect.ValueOf(fn)
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

	return func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Dump()
		params := reflect.New(fnType.In(1))
		if err := s.webApiDecode(&ctx, req, params.Interface()); err != nil {
			s.webApiHandleResp(&ctx, w, nil, err)
			return
		}

		vals := fnVal.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			params.Elem(),
		})
		err, ok := vals[1].Interface().(error)
		if !ok {
			err = nil
		}
		s.webApiHandleResp(&ctx, w, vals[0].Interface(), err)
	}
}

func (s *WebServer) webApiDecode(ctx *context.T, req *http.Request, arg interface{}) error {
	// 创建并初始化context
	h := ctx.InitRequestHeap(nil)
	h.Start = time.Now()
	ctx.SetRequest(req)

	if s.IsDownload && req.Method == http.MethodGet {
		// Get 的下载没有body
		data := req.URL.Query().Get("_")
		if err := json.Unmarshal([]byte(data), arg); err != nil {
			return err
		}
	} else {
		requestBody, err := ioutil.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return err
		}
		err = json.Unmarshal(requestBody, arg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *WebServer) webApiHandleResp(ctx *context.T, w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		errorEncoder(err, w)
		return
	} else {
		if c := ctx.GetResponseCookie(); c != nil {
			for _, ck := range c {
				http.SetCookie(w, ck)
			}
		}
	}

	if s.IsDownload && err == nil {
		file := resp.(Filer)
		filename := file.GetFilename()
		if filename == "" {
			panic("filename must be set")
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		http.ServeContent(w, ctx.GetRequest(), filename, time.Now(), bytes.NewReader(file.GetData()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		ret, _ := json.Marshal(resp)
		w.Write(ret)
	}
}

func (s *WebServer) AutoRouter(c interface{}) {
	if s.mux == nil {
		s.mux = mux.NewRouter()
	}

	pathPrefix := fmt.Sprintf("/api/%s.%s", mconf.GetPkgName(), getServiceName(c))
	subRouter := s.mux.PathPrefix(pathPrefix).Subrouter()
	vf := reflect.ValueOf(c)
	ctx := context.Dump()
	for i := 0; i < vf.NumMethod(); i++ {
		func(i int) {
			subRouter.HandleFunc(fmt.Sprintf("/%v", vf.Type().Method(i).Name), generateHandler(ctx, vf.Method(i))).Methods("POST")
		}(i)
	}
}

func (s *WebServer) Close() error {
	err := s.srv.Shutdown(nil)
	grace.DoneTask()
	return errors.Trace(err)
}

func getServiceName(s interface{}) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	name := path.Base(t.PkgPath())
	first := strings.ToUpper(string(name[0]))
	return first + name[1:]
}

func errorEncoder(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
	log.NewEx(-1).Error(err)
}

type errorWrapper struct {
	Error string `json:"error"`
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
