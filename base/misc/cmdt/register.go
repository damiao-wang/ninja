package cmdt

import (
	"log"
	"ninja/base/mconf"
	"ninja/base/misc/stack"
	"os"
	"path"
	"reflect"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
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

func init() {
	mconf.ParseFlag(RootCmd.PersistentFlags())
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
	Run(ctx context.Context) error
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

func registerServiceEx(name, desc string, s Servicer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: desc,
		Run: func(c *cobra.Command, args []string) {
			defer s.Close()
			if err := registerServicer(s); err != nil {
				log.Println(err)
				return
			}
			s.Run(context.Background())
			return
		},
	}
	return cmd
}
