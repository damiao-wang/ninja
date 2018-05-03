package grace

import (
	"errors"
	"sync"
	"sync/atomic"

	"ninja/base/misc/log"
)

var (
	isClosed   int32
	processing sync.WaitGroup
	ErrClosing = errors.New("server is closing")
)

func ExitMark() {
	atomic.StoreInt32(&isClosed, 1)
	log.Info("ready exit. wait processing task done.")
}

func WaitTaskDone() {
	processing.Wait()
	log.Info("all processing task done. exit")
}

func CountInTask() bool {
	if atomic.LoadInt32(&isClosed) == 1 {
		return false
	}
	processing.Add(1)
	return true
}

func DoneTask() {
	processing.Done()
}
