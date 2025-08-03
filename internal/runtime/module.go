package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

func (r *JSRuntime) setupRequire() {
	r.vm.Set("require", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(r.vm.ToValue("require() needs a path"))
		}

		modPath := call.Arguments[0].String()

		absPath, err := filepath.Abs(modPath)
		if err != nil {
			panic(r.vm.ToValue("Invalid path: " + err.Error()))
		}
		if mod, ok := r.moduleCache[absPath]; ok && mod.loaded {
			return mod.exports
		}

		code, err := os.ReadFile(absPath)
		if err != nil {
			panic(r.vm.ToValue("Cannot read module: " + err.Error()))
		}

		codeStr := strings.TrimPrefix(string(code), "\uFEFF")

		wrapped := fmt.Sprintf(`(function(exports, require, module, __filename, __dirname) {%s})`, codeStr)

		ast, err := parser.ParseFile(nil, absPath, wrapped, 0)
		if err != nil {
			panic(r.vm.ToValue("Failed to parse module: " + err.Error()))
		}

		program, err := goja.CompileAST(ast, false)
		if err != nil {
			panic(r.vm.ToValue("Failed to compile module: " + err.Error()))
		}

		fnVal, err := r.vm.RunProgram(program)
		if err != nil {
			panic(r.vm.ToValue("Module execution failed: " + err.Error()))
		}

		fn, ok := goja.AssertFunction(fnVal)
		if !ok {
			panic("Invalid module function")
		}

		exports := r.vm.NewObject()
		module := r.vm.NewObject()
		_ = module.Set("exports", exports)

		r.moduleCache[absPath] = &Module{
			exports: exports,
			loaded:  false,
		}

		dir := filepath.Dir(absPath)
		_, err = fn(goja.Undefined(), exports, r.vm.Get("require"), module, r.vm.ToValue(absPath), r.vm.ToValue(dir))
		if err != nil {
			panic(r.vm.ToValue("Module execution error: " + err.Error()))
		}

		// Store final exports
		finalExports := module.Get("exports")
		r.moduleCache[absPath].exports = finalExports
		r.moduleCache[absPath].loaded = true

		return finalExports
	})
}
