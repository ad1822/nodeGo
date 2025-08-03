package runtime

import (
	"github.com/dop251/goja"
)

func (r *JSRuntime) initImmediate() {
	r.vm.Set("setImmediate", func(call goja.FunctionCall) goja.Value {
		cb, ok := goja.AssertFunction(call.Argument(0))
		if !ok {
			panic(r.vm.ToValue("First argument must be a function"))
		}

		r.immediateQueue <- func() {
			_, err := cb(goja.Undefined(), nil)
			if err != nil {
				println("[setImmediate] callback error:", err.Error())
			}
		}

		return goja.Undefined()
	})
}
