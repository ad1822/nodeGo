# Event Emitter

- An EventEmitter is a public-subscriber (pub/ sub) mechanism in code
- **Publisher** (via `emit`) announces events
- **Subscriber** (via `on`) registers functions (listeners) to react when whose events happen
- Unsubscriber (via `off`) removes those listeners

## Core Methods

## 1. on(event, listener)

- Purpose: Register a listener for an event.

- Behavior:

  - Look up the event name in your internal storage (often map[string][]func(...)). If no listeners yet, create a new list. Append the new listener.

  - Example:

  ```js
  emitter.on('data', (chunk) => {
    console.log('Got chunk:', chunk);
  });
  ```

2.  emit(event, ...args)

- Purpose: Trigger all listeners for the given event, passing arguments.

- Behavior:

  - Look up all listeners for that event. Call each listener in the order they were registered. Pass ...args to them.

  - Example:

  ```js
  emitter.emit('data', 'Hello World');
  // All "data" listeners get called with "Hello World"
  ```

3.  off(event, listener)

- Purpose: Remove a specific listener from an event.

- Behavior:

  - Look up the listeners for that event. Remove any that match the given listener reference.

  - Example:

  ```js
  const handler = (msg) => console.log(msg);
  emitter.on('msg', handler);
  emitter.off('msg', handler);
  ```
