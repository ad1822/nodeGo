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
	immediateQueue chan func()
	nextTickQueue  chan func()
	nextTimerId    int
	activeTimers   map[int]bool
}

// Have to exist explicitly
func (r *JSRuntime) RunEventLoop() {
	for {
		select {
		case nextTick := <-r.nextTickQueue:
			nextTick()

			for len(r.nextTickQueue) > 0 {
				nextTick := <-r.nextTickQueue
				nextTick()
			}

		case micro := <-r.microTaskQueue:
			micro()

			for len(r.microTaskQueue) > 0 {
				micro := <-r.microTaskQueue
				micro()
			}

		case immediate := <-r.immediateQueue:
			immediate()

		case task := <-r.taskQueue:
			task()

		default:
			r.vm.RunString("")
			time.Sleep(10 * time.Millisecond)

		}
	}

}

func New() *JSRuntime {
	r := &JSRuntime{
		vm:             goja.New(),
		microTaskQueue: make(chan func(), 100), // Micro Task Queue for callbacks and Promises
		immediateQueue: make(chan func(), 100), // Immediate Queue for callbacks before timers
		taskQueue:      make(chan func(), 100), // buffered task queue
		nextTickQueue:  make(chan func(), 100), // Next tick Queue for callbacks before microtask queue
		nextTimerId:    0,
		activeTimers:   make(map[int]bool),
	}
	r.initConsole()
	r.initNextTickQueue()
	r.initImmediate()
	r.initMicroTaskQueue()
	r.initTimers()
	// r.initClearTimers()
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
