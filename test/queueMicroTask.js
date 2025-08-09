const log = require('./test/logger.js');

queueMicrotask(() => {
  log('microtask 1');
});

Promise.resolve().then(() => {
  console.error('Promise resolved');
});
