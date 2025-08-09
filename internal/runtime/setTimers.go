package runtime

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
)

// Set timers (setTimeout, setInterval) APIS
func (r *JSRuntime) initTimers() {
	r.vm.Set("setTimeout", func(call goja.FunctionCall) goja.Value {
		r.nextTimerId++
		id := r.nextTimerId

		r.activeTimers[id] = true

		// fmt.Println("ID of setTimeout :", id)
		if len(call.Arguments) < 2 {
			panic(r.vm.ToValue("setTimeout expects a callback and delay"))
		}

		obj := call.Arguments[0].ToObject(r.vm)

		cb, ok := goja.AssertFunction(obj)

		if !ok {
			panic(r.vm.ToValue("First argument must be a function"))
		}

		delay := call.Arguments[1].ToInteger()

		go func() {
			// Set delay for that particular time
			time.Sleep(time.Duration(delay) * time.Millisecond)

			// Append that function in taskQueue
			r.taskQueue <- func() {
				_, err := cb(goja.Undefined())
				if err != nil {
					fmt.Println("[Go] Error in callback:", err)
				}
			}

		}()

		return r.vm.ToValue(id)

	})

	// setInterval API
	r.vm.Set("setInterval", func(call goja.FunctionCall) goja.Value {
		r.nextTimerId++
		id := r.nextTimerId

		r.activeTimers[id] = true

		// fmt.Println("ID of setInterval :", id)
		if len(call.Arguments) < 2 {
			panic(r.vm.ToValue("setTimeout expects a callback and delay"))
		}

		cb := call.Argument(0)
		delay := call.Argument(1).ToInteger()

		if callable, ok := goja.AssertFunction(cb); ok {
			go func() {
				ticker := time.NewTicker(time.Duration(delay) * time.Millisecond)
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						r.taskQueue <- func() {
							_, err := callable(goja.Undefined())
							if err != nil {
								fmt.Println("setInterval error:", err)
							}
						}
						// case <-r.done: // Optional shutdown support
						// return
					}
				}
			}()
		}

		return r.vm.ToValue(id)
	})

}
