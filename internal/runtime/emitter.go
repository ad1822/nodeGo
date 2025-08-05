package runtime

import "github.com/dop251/goja"

func NewEventEmitter(vm *goja.Runtime) *EventEmitter {
	return &EventEmitter{
		vm:     vm,
		events: make(map[string][]goja.Callable),
	}
}

func (e *EventEmitter) On(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		return goja.Undefined()
	}

	event := call.Argument(0).String()
	fn, ok := goja.AssertFunction(call.Argument(1))
	if !ok {
		panic(e.vm.ToValue("Event listener must be a function"))
	}

	e.events[event] = append(e.events[event], fn)
	return goja.Undefined()
}

func (e *EventEmitter) Emit(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		return goja.Undefined()
	}

	event := call.Argument(0).String()
	args := call.Arguments[1:]

	listeners, ok := e.events[event]
	if !ok {
		return goja.Undefined()
	}

	for _, fn := range listeners {
		_, err := fn(goja.Undefined(), args...)
		if err != nil {
			panic(err)
		}
	}

	return goja.Undefined()
}

// func (e *EventEmitter) Off(call goja.FunctionCall) goja.Value {
// 	if len(call.Arguments) < 2 {
// 		return goja.Undefined()
// 	}

// 	event := call.Argument(0).String()
// 	targetFn, ok := goja.AssertFunction(call.Argument(1))
// 	if !ok {
// 		return goja.Undefined()
// 	}

// 	current := e.events[event]
// 	filtered := []goja.Callable{}

// 	for _, fn := range current {
// 		if !e.vm.StrictEquals(fn, targetFn) {
// 		}
// 	}

// 	e.events[event] = filtered
// 	return goja.Undefined()
// }

func (r *JSRuntime) initEventEmitter() {
	// class := r.vm.NewObject()

	emitterConstructor := func(call goja.ConstructorCall) *goja.Object {
		emitter := NewEventEmitter(r.vm)

		obj := call.This

		obj.Set("on", emitter.On)
		obj.Set("emit", emitter.Emit)
		// obj.Set("off", emitter.Off) // optional

		return obj
	}

	r.vm.Set("EventEmitter", emitterConstructor)
}
