package runtime

import (
	"fmt"

	"github.com/dop251/goja"
)

func (r *JSRuntime) initMicroTaskQueue() {
	r.vm.Set("queueMicrotask", func(call goja.FunctionCall) goja.Value {
		cb, ok := goja.AssertFunction(call.Argument(0))
		if !ok {
			panic("Not a function")
		}

		r.microTaskQueue <- func() {
			_, err := cb(goja.Undefined())
			if err != nil {
				fmt.Println("[Go] Error in microtask callback:", err)
			}
		}

		return goja.Undefined()
	})
}
