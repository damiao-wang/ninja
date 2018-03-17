package cmdt

import (
	"reflect"
	"path"
	"fmt"
	"net/http"
	"time"
	"ninja/base/misc/stack"
	"github.com/spf13/cobra"
	"github.com/gorilla/mux"
	"ninja/base/misc/router"
	"github.com/urfave/negroni"
	"log"
	"os/signal"
	"os"
	"syscall"
)

var RootCmd = &cobra.Command{
	Use: stack.GetRootService(),
}

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "服务列表",
}

func SetName(name string) {
	RootCmd.Short = name
}

func Execute() error {
	// 用于运维创建索引
	if os.Getenv("MGOINDEX") != "" {
		return nil
	}

	return RootCmd.Execute()
}

func GetServiceName(s interface{}) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return path.Base(t.PkgPath())
}

type Servicer interface {
	Register() error
	Desc() string
	// Run(ctx context.Context) error
	Close() error
}

type Initer interface {
	Init(srv interface{}, pkg string, register func() error) error
}

func RegisterService(s Servicer) *cobra.Command {
	name := GetServiceName(s)
	cmd := registerServiceEx(name, s.Desc(), s)
	if cmd.Use == "main" {
		RootCmd.Run = cmd.Run
		return RootCmd
	}

	ServiceCmd.AddCommand(cmd)
	addCommandOnce(RootCmd, ServiceCmd)
	return cmd
}

func addCommandOnce(parent, me *cobra.Command) {
	for _, sub := range parent.Commands() {
		if sub == me {
			return
		}
	}
	parent.AddCommand(me)
}

func registerServicer(s Servicer) error {
	name := GetServiceName(s)
	if in, ok := s.(Initer); ok {
		return in.Init(s, name, s.Register)
	}
	return s.Register()
}

func registerServiceEx(name, desc string, s interface{}) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: desc,
		Run: func(c *cobra.Command, args []string) {
			if err := registerServicer(s.(Servicer)); err != nil {
				log.Println(err)
				return
			}
			r := mux.NewRouter().Path(fmt.Sprintf("/api/name")).Subrouter()
			router.AutoRouter(r, s)
			n := negroni.Classic()
			n.UseHandler(r)
			
			srv := &http.Server{
				Addr: ":8080",
				Handler: n,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}
		
			go func() {
				log.Println(srv.ListenAndServe())
				log.Println("server shutdown")
			}()
		
			// Handle SIGINT and SIGTERM.
			ch := make(chan os.Signal)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			log.Println(<-ch)
		
			// Stop the service gracefully.
			log.Println(srv.Shutdown(nil))
			s.(Servicer).Close()
		
			// Wait gorotine print shutdown message
			time.Sleep(time.Second * 5)
			log.Println("done.")
			return
		},
	}
	return cmd
}
