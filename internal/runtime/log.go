package runtime

import (
	"fmt"

	"github.com/dop251/goja"
)

func (r *JSRuntime) initConsole() {
	console := map[string]func(goja.FunctionCall) goja.Value{
		"log": func(call goja.FunctionCall) goja.Value {
			for _, arg := range call.Arguments {
				fmt.Print(arg.Export(), " ")
			}
			fmt.Println()
			return goja.Undefined()
		},
	}
	r.vm.Set("console", console)
}
