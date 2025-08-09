const emitter = new EventEmitter();

function listener1(data) {
  console.log('listener1 received:', data);
}

function listener2(data) {
  console.log('listener2 received:', data);
}

// Register listeners
emitter.on('test', listener1);
emitter.on('test', listener2);

console.log('---- First emit ----');
emitter.emit('test', 'Hello World');
// Expected:
// listener1 received: Hello World
// listener2 received: Hello World

// Remove listener1
emitter.off('test', listener1);

console.log('---- Second emit (after removing listener1) ----');
emitter.emit('test', 'Hello Again');
// Expected:
// listener2 received: Hello Again

// Remove listener2
emitter.off('test', listener2);

console.log('---- Third emit (no listeners) ----');
emitter.emit('test', 'No one hears this');
// Expected:
// (nothing printed)
