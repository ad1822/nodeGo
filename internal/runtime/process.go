package runtime

import (
	"os"
	"strings"

	"github.com/dop251/goja"
)

func (r *JSRuntime) initProcess() {
	process := r.vm.NewObject()

	process.Set("cwd", func(goja.FunctionCall) goja.Value {
		dir, err := os.Getwd()
		if err != nil {
			return r.vm.ToValue("")
		}
		return r.vm.ToValue(dir)
	})

	argv := []string{}
	for _, arg := range os.Args {
		argv = append(argv, arg+"\n")
	}
	process.Set("argv", r.vm.ToValue(argv))

	envObj := r.vm.NewObject()
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envObj.Set(parts[0], parts[1]+"\n")
		}
	}
	process.Set("env", envObj)

	process.Set("exit", func(call goja.FunctionCall) goja.Value {
		code := 0
		if len(call.Arguments) > 0 {
			code = int(call.Argument(0).ToInteger())
		}
		os.Exit(code)
		return goja.Undefined()
	})

	r.vm.Set("process", process)
}
