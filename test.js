const startTime = Date.now();

function log(msg) {
  const elapsed = ((Date.now() - startTime) / 1000).toFixed(2);
  console.log(`[+${elapsed}s] ${msg}`);
}

log('start');

i = 0;

Promise.resolve().then(() => {
  log('Promise resolved');
});

setInterval(() => {
  if (i == 5) return;
  i++;
  log(`Interval ${i}`);
}, 2000);

queueMicrotask(() => {
  log('microtask 1');
});

setTimeout(() => {
  log('setTimeout 2');

  queueMicrotask(() => {
    log('microtask 2 (inside setTimeout 2)');
  });
}, 2000);

log('end');
