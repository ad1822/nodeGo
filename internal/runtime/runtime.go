package runtime

import (
	"os"
	"time"

	"github.com/dop251/goja"
)

type JSRuntime struct {
	vm             *goja.Runtime
	taskQueue      chan func() // Chan is for channel
	microTaskQueue chan func()
}

// Have to exist explicitly
func (r *JSRuntime) RunEventLoop() {
	for {
		select {
		case micro := <-r.microTaskQueue:
			micro()

			for len(r.microTaskQueue) > 0 {
				micro := <-r.microTaskQueue
				micro()
			}

		case task := <-r.taskQueue:
			task()

		default:
			time.Sleep(10 * time.Millisecond)

		}
	}

}

func New() *JSRuntime {
	r := &JSRuntime{
		vm:             goja.New(),
		taskQueue:      make(chan func(), 100), // buffered task queue
		microTaskQueue: make(chan func(), 100),
	}
	r.initConsole()
	r.initTimers()
	r.initMicroTaskQueue()
	return r
}

func (r *JSRuntime) RunScript(path string) error {
	code, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = r.vm.RunString(string(code))
	if err != nil {
		return err
	}

	// Run your event loop to wait for async tasks
	r.RunEventLoop()
	return nil
}
