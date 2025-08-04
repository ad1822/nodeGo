package runtime

import (
	"os"
	"time"

	"github.com/dop251/goja"
)

type Module struct {
	exports goja.Value
	loaded  bool
}

type JSRuntime struct {
	vm             *goja.Runtime
	taskQueue      chan func() // Chan is for channel
	microTaskQueue chan func()
	immediateQueue chan func()
	nextTickQueue  chan func()
	nextTimerId    int
	activeTimers   map[int]bool
	moduleCache    map[string]*Module
}

// Have to exist explicitly
func (r *JSRuntime) RunEventLoop() {
	for {
		// Drain nextTick queue
		for {
			select {
			case fn := <-r.nextTickQueue:
				fn()
			default:
				goto drainMicrotasks
			}
		}

	drainMicrotasks:
		for {
			select {
			case fn := <-r.microTaskQueue:
				fn()
			default:
				goto drainImmediate
			}
		}

	drainImmediate:
		for {
			select {
			case fn := <-r.immediateQueue:
				fn()
			default:
				goto drainTimers
			}
		}

	drainTimers:
		for {
			select {
			case fn := <-r.taskQueue:
				fn()
			default:
				goto idle
			}
		}

	idle:
		// If all queues are empty, wait briefly
		time.Sleep(1 * time.Millisecond)
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
		moduleCache:    make(map[string]*Module),
	}
	r.initConsole()
	r.initNextTickQueue()
	r.initImmediate()
	r.initMicroTaskQueue()
	r.initTimers()
	r.setupRequire()
	r.initProcess()
	// r.setupFSModule()
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
