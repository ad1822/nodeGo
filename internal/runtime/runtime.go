package runtime

import (
	"os"

	"github.com/dop251/goja"
)

type JSRuntime struct {
	vm        *goja.Runtime
	taskQueue chan func() // Chan is for channel
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
