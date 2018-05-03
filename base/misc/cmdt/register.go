package cmdt

import (
	"os"
	"os/signal"
	"path"
	"reflect"
	"syscall"

	"ninja/base/mconf"
	"ninja/base/misc/context"
	"ninja/base/misc/grace"
	"ninja/base/misc/stack"
	"ninja/base/trace/sentry"

	"github.com/spf13/cobra"
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
	Run(ctx context.T) error
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
			ctx, cf := context.WithCancel(context.Dump())
			if err := registerServicer(s); err != nil {
				ctx.LogErrorEx(1, err)
				return
			}

			// 捕获退出信号，gracefull退出
			catchSignal(cf)

			var err error
			sentry.CapturePanicEx(0, "", func() {
				err = s.Run(ctx)
				s.Close()
				grace.WaitTaskDone()
				if err != nil {
					ctx.LogErrorEx(1, err)
					os.Exit(1)
				}
			}, nil)

			return
		},
	}
	return cmd
}

// cf 是在响应退出信号时必须要执行的
func catchSignal(cf func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)
	go func() {
		var exitSignalCount int
		for {
			select {
			case s := <-signalChan:
				switch s {
				case os.Interrupt, syscall.SIGTERM, syscall.SIGHUP:
					exitSignalCount++
					if exitSignalCount > 1 {
						context.LogInfo("force exit")
						os.Exit(1)
					}
					grace.ExitMark()
					cf()
				default:

				}
			}
		}
	}()

}
