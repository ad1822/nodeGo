package runtime

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
)

func (r *JSRuntime) initTimers() {
	r.vm.Set("setTimeout", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(r.vm.ToValue("setTimeout expects a callback and delay"))
		}

		// fmt.Println("Inside Go")

		obj := call.Arguments[0].ToObject(r.vm)
		// fmt.Println(obj)
		cb, ok := goja.AssertFunction(obj)
		if !ok {
			panic(r.vm.ToValue("First argument must be a function"))
		}

		delay := call.Arguments[1].ToInteger()
		// fmt.Println(delay)

		go func() {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			fmt.Println("[Go] Timeout finished, pushing callback to taskQueue")

			// fmt.Println("Before pushing task, queue length:", len(r.taskQueue))
			r.taskQueue <- func() {
				fmt.Println("[Go] Executing setTimeout callback from taskQueue")
				_, err := cb(goja.Undefined())
				if err != nil {
					fmt.Println("[Go] Error in callback:", err)
				}
			}
			// fmt.Println("After pushing task, queue length:", len(r.taskQueue))

			fmt.Println("[Go] Callback pushed successfully")
		}()

		return goja.Undefined()
	})
}

// Have to exist explicitly
func (r *JSRuntime) RunEventLoop() {
	for {
		select {
		case task := <-r.taskQueue:
			// fmt.Println("Queue length:", len(r.taskQueue))
			task()
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
