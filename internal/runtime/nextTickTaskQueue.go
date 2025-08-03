package runtime

import (
	"fmt"

	"github.com/dop251/goja"
)

func (r *JSRuntime) initNextTickQueue() {
	r.vm.Set("process", map[string]any{})

	process := r.vm.Get("process").ToObject(r.vm)
	process.Set("nextTick", func(call goja.FunctionCall) goja.Value {
		cb, ok := goja.AssertFunction(call.Argument(0))
		if !ok {
			panic(r.vm.ToValue("First argument must be a function"))
		}

		r.nextTickQueue <- func() {
			_, err := cb(goja.Undefined())
			// fmt.Println("[GO NEXT TICK] GOES")
			if err != nil {
				fmt.Println("[Go] Error in nextTick:", err)
			}
		}

		return goja.Undefined()
	})

}
