const startTime = Date.now();

function log(msg) {
  const elapsed = ((Date.now() - startTime) / 1000).toFixed(2);
  console.log(`[+${elapsed}s] ${msg}`);
}

// log('start');

i = 0;

setImmediate(() => console.warn('immediate'));

Promise.resolve().then(() => {
  console.error('Promise resolved');
});

// setInterval(() => {
//   if (i == 5) return;
//   i++;
//   console.warn(`Interval ${i}`);
// }, 2000);

queueMicrotask(() => {
  log('microtask 1');
});

process.nextTick(() => {
  console.warn('NEXT TICK 1');
});

let id = setInterval(() => {
  log('Interval');
}, 1000);

process.nextTick(() => {
  console.warn('NEXT TICK 2');
});

// let id2 = setTimeout(() => {
//   log('Timeout');
// }, 2000);

// log('end');
