- Go rountines and Channnles
- How to stop go routines
- Clean timeout and interval

### Order of Execution

```
Sync (Promises) -> NextTick Queue -> MicroTask Queue -> Immediate Queue -> Timers (setInterval, setTimeout)


[ JS Engine ]
|
|— execute script
|   |— enqueue all Promises immediately (inside Goja)
|
|— return control to Go
|
[ Go Event Loop ]
|
|— drain nextTickQueue
|— drain microTaskQueue
|— drain immediateQueue
|— drain timerQueue

```

- ✅ `console` (log, warn, error)
- ✅ `setTimeout`, `setInterval`,
- [ ] `clearTimeout`, `clearInterval`
- ✅ `setImmediate`
- ✅ Microtask queue (`Promise.then`, `queueMicrotask`)
- ✅ Event loop (custom scheduler)

---

### 🚀 **Next Features to Implement**

#### 1. **`process.nextTick()`**

- Executes _before_ other microtasks.
- Goes into a separate queue with higher priority than `Promise.then`.

**Order becomes:**

```text
1. sync
2. nextTick queue
3. microTask queue
4. setImmediate
5. timers
```

> This is critical if you want Node.js-level behavior.

---

#### 2. **`require()` or Module Support**

- Start with basic file-based `require` that loads `.js` files using Go’s file system.
- Use `vm.RunString()` with file content.
- Store module exports in a map to avoid duplicate executions (simple cache).

---

#### 3. **Basic `fs` Module (as API mock)**

- Implement a very limited version of Node’s `fs` module.
- Expose methods like:

  - `fs.readFile(path, callback)`
  - `fs.writeFile(path, data, callback)`

You can wire this using Go’s `os.ReadFile` / `os.WriteFile` and send callback into the `taskQueue`.

---

#### 4. **EventEmitter Class**

- Implement `on`, `emit`, `off`
- You’ll need to track events and associated callbacks in Go.
- Useful to simulate native Node.js behavior.

---

#### 5. **Uncaught Error Handler**

- Trap unhandled errors from promises or sync code
- Provide a way to register a listener like:

```js
process.on('uncaughtException', (err) => { ... });
```

---

#### 6. **Built-in `process` Object**

- At least define:

  - `process.nextTick`
  - `process.env`
  - `process.pid`, etc.

---

### 📌 Priority Suggestion

| Step | Feature               | Why It Matters                     |
| ---- | --------------------- | ---------------------------------- |
| 1    | `process.nextTick()`  | For completeness of event loop     |
| 2    | `require()` + modules | Enables modular JS                 |
| 3    | `fs` module           | File I/O with async pattern        |
| 4    | `EventEmitter`        | Needed for `fs`, server simulation |
| 5    | Error Handling        | Stability and developer UX         |
| 6    | `process` object      | Completes Node.js-like environment |

---
