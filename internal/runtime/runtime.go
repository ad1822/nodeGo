package runtime

import (
	"os"
	"time"

	"github.com/dop251/goja"
)

type JSRuntime struct {
	vm        *goja.Runtime
	taskQueue chan func() // Chan is for channel
}

// Have to exist explicitly
func (r *JSRuntime) RunEventLoop() {
	for {
		select {
		case task := <-r.taskQueue:
			// fmt.Println("Queue length:", len(r.taskQueue))
			task()
			// r.vm.RunString("")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func New() *JSRuntime {
	r := &JSRuntime{
		vm:        goja.New(),
		taskQueue: make(chan func(), 100), // buffered task queue
	}
	r.initConsole()
	r.initTimers()

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
