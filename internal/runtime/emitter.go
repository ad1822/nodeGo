package runtime

import (
	"fmt"

	"github.com/dop251/goja"
)

func NewEventEmitter(vm *goja.Runtime) *EventEmitter {
	return &EventEmitter{
		vm:     vm,
		events: make(map[string][]goja.Value),
	}
}

func (e *EventEmitter) On(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		return goja.Undefined()
	}

	event := call.Argument(0).String()
	fn := call.Argument(1)
	if _, ok := goja.AssertFunction(fn); !ok {
		return goja.Undefined()
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
		cb, _ := goja.AssertFunction(fn)
		_, err := cb(goja.Undefined(), args...)
		if err != nil {
			fmt.Println("[Go] listener error:", err)
		}
	}

	return goja.Undefined()
}

func (e *EventEmitter) Off(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 2 {
		return goja.Undefined()
	}

	event := call.Argument(0).String()
	target := call.Argument(1)
	listeners := e.events[event]
	for i, fn := range listeners {
		if fn.Equals(target) { // JS === comparison
			listeners = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
	e.events[event] = listeners
	return goja.Undefined()

}

func (r *JSRuntime) initEventEmitter() {
	// class := r.vm.NewObject()

	emitterConstructor := func(call goja.ConstructorCall) *goja.Object {
		emitter := NewEventEmitter(r.vm)

		obj := call.This

		obj.Set("on", emitter.On)
		obj.Set("emit", emitter.Emit)
		obj.Set("off", emitter.Off) // optional

		return obj
	}

	r.vm.Set("EventEmitter", emitterConstructor)
}
