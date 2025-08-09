package runtime

import (
	"fmt"
	"io"
	"net"

	"github.com/dop251/goja"
)

// ======================
// EventEmitter
// ======================

// type EventEmitter struct {
// 	vm     *goja.Runtime
// 	events map[string][]goja.Value
// }

// func NewEventEmitter(vm *goja.Runtime) *EventEmitter {
// 	return &EventEmitter{
// 		vm:     vm,
// 		events: make(map[string][]goja.Value),
// 	}
// }

func (e *EventEmitter) EmitGo(event string, args ...goja.Value) {
	if list, ok := e.events[event]; ok {
		for _, fn := range list {
			if cb, ok := goja.AssertFunction(fn); ok {
				_, err := cb(goja.Undefined(), args...)
				if err != nil {
					fmt.Println("[Go] listener error:", err)
				}
			}
		}
	}
}

func (e *EventEmitter) OnGo(event string, fn goja.Value) {
	if _, ok := goja.AssertFunction(fn); !ok {
		return // ignore non-functions
	}
	e.events[event] = append(e.events[event], fn)
}

func (e *EventEmitter) OffGo(event string, target goja.Value) {
	if listeners, ok := e.events[event]; ok {
		for i, fn := range listeners {
			if fn.Equals(target) {
				e.events[event] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

// ======================
// TCP Socket
// ======================

type Socket struct {
	vm      *goja.Runtime
	conn    net.Conn
	emitter *EventEmitter
}

func NewSocket(vm *goja.Runtime, conn net.Conn) *Socket {
	s := &Socket{
		vm:      vm,
		conn:    conn,
		emitter: NewEventEmitter(vm),
	}
	go s.readLoop()
	return s
}

func (s *Socket) readLoop() {
	buf := make([]byte, 1024)
	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				s.emitter.EmitGo("error", s.vm.ToValue(err.Error()))
			}
			s.emitter.EmitGo("end")
			s.emitter.EmitGo("close")
			return
		}
		s.emitter.EmitGo("data", s.vm.ToValue(string(buf[:n])))
	}
}

func (s *Socket) Write(data string) {
	_, err := s.conn.Write([]byte(data))
	if err != nil {
		s.emitter.EmitGo("error", s.vm.ToValue(err.Error()))
	}
}

func (s *Socket) End() {
	s.conn.Close()
	s.emitter.EmitGo("end")
	s.emitter.EmitGo("close")
}

func (s *Socket) ToJS() *goja.Object {
	obj := s.vm.NewObject()

	obj.Set("write", func(call goja.FunctionCall) goja.Value {
		s.Write(call.Argument(0).String())
		return goja.Undefined()
	})

	obj.Set("end", func(goja.FunctionCall) goja.Value {
		s.End()
		return goja.Undefined()
	})

	obj.Set("on", func(call goja.FunctionCall) goja.Value {
		s.emitter.OnGo(call.Argument(0).String(), call.Argument(1))
		return goja.Undefined()
	})

	obj.Set("off", func(call goja.FunctionCall) goja.Value {
		s.emitter.OffGo(call.Argument(0).String(), call.Argument(1))
		return goja.Undefined()
	})

	return obj
}

// ======================
// TCP Server
// ======================

type TCPServer struct {
	runtime  *JSRuntime
	emitter  *EventEmitter
	listener net.Listener
}

func (srv *TCPServer) Listen(port int, host string) {
	addr := fmt.Sprintf("%s:%d", host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		srv.emitter.EmitGo("error", srv.runtime.vm.ToValue(err.Error()))
		return
	}
	srv.listener = ln
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				srv.emitter.EmitGo("error", srv.runtime.vm.ToValue(err.Error()))
				return
			}
			socket := NewSocket(srv.runtime.vm, conn)
			srv.emitter.EmitGo("connection", socket.ToJS())
		}
	}()
}

func (srv *TCPServer) Close() {
	if srv.listener != nil {
		srv.listener.Close()
		srv.emitter.EmitGo("close")
	}
}

func (srv *TCPServer) ToJS() *goja.Object {
	obj := srv.runtime.vm.NewObject()

	obj.Set("listen", func(call goja.FunctionCall) goja.Value {
		port := int(call.Argument(0).ToInteger())
		host := "0.0.0.0"
		if len(call.Arguments) > 1 {
			host = call.Argument(1).String()
		}
		srv.Listen(port, host)
		return goja.Undefined()
	})

	obj.Set("close", func(goja.FunctionCall) goja.Value {
		srv.Close()
		return goja.Undefined()
	})

	obj.Set("on", func(call goja.FunctionCall) goja.Value {
		srv.emitter.OnGo(call.Argument(0).String(), call.Argument(1))
		return goja.Undefined()
	})

	obj.Set("off", func(call goja.FunctionCall) goja.Value {
		srv.emitter.OffGo(call.Argument(0).String(), call.Argument(1))
		return goja.Undefined()
	})

	return obj
}

// ======================
// TCP Client
// ======================

type TCPClient struct {
	runtime *JSRuntime
	emitter *EventEmitter
	conn    net.Conn
}

func (c *TCPClient) connect(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		c.emitter.EmitGo("error", c.runtime.vm.ToValue(err.Error()))
		return
	}
	c.conn = conn
	c.emitter.EmitGo("connect")
	go c.readLoop()
}

func (c *TCPClient) readLoop() {
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				c.emitter.EmitGo("error", c.runtime.vm.ToValue(err.Error()))
			}
			c.emitter.EmitGo("end")
			c.emitter.EmitGo("close")
			return
		}
		c.emitter.EmitGo("data", c.runtime.vm.ToValue(string(buf[:n])))
	}
}

func (c *TCPClient) Write(data string) {
	if c.conn != nil {
		if _, err := c.conn.Write([]byte(data)); err != nil {
			c.emitter.EmitGo("error", c.runtime.vm.ToValue(err.Error()))
		}
	}
}

func (c *TCPClient) End() {
	if c.conn != nil {
		c.conn.Close()
		c.emitter.EmitGo("end")
		c.emitter.EmitGo("close")
	}
}

func (c *TCPClient) ToJS() *goja.Object {
	obj := c.runtime.vm.NewObject()

	obj.Set("write", func(call goja.FunctionCall) goja.Value {
		c.Write(call.Argument(0).String())
		return goja.Undefined()
	})

	obj.Set("end", func(goja.FunctionCall) goja.Value {
		c.End()
		return goja.Undefined()
	})

	obj.Set("on", func(call goja.FunctionCall) goja.Value {
		c.emitter.OnGo(call.Argument(0).String(), call.Argument(1))
		return goja.Undefined()
	})

	obj.Set("off", func(call goja.FunctionCall) goja.Value {
		c.emitter.OffGo(call.Argument(0).String(), call.Argument(1))
		return goja.Undefined()
	})

	return obj
}

// ======================
// initNet binding
// ======================

func (r *JSRuntime) initNet() {
	netObj := r.vm.NewObject()

	netObj.Set("createServer", func(call goja.FunctionCall) goja.Value {
		var listener goja.Value
		if len(call.Arguments) > 0 {
			listener = call.Argument(0)
		}
		server := &TCPServer{
			runtime: r,
			emitter: NewEventEmitter(r.vm),
		}
		if listener != nil {
			server.emitter.OnGo("connection", listener)
		}
		return server.ToJS()
	})

	netObj.Set("createConnection", func(call goja.FunctionCall) goja.Value {
		options := call.Argument(0).ToObject(r.vm)
		host := "127.0.0.1"
		port := 0
		if p := options.Get("port"); p != nil {
			port = int(p.ToInteger())
		}
		if h := options.Get("host"); h != nil {
			host = h.String()
		}

		var connectListener goja.Value
		if len(call.Arguments) > 1 {
			connectListener = call.Argument(1)
		}

		client := &TCPClient{
			runtime: r,
			emitter: NewEventEmitter(r.vm),
		}
		if connectListener != nil {
			client.emitter.OnGo("connect", connectListener)
		}
		go client.connect(host, port)
		return client.ToJS()
	})

	r.vm.Set("net", netObj)
}
