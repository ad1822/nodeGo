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

let id = setInterval(() => {
  log('Interval');
}, 1000);

let id2 = setTimeout(() => {
  log('Timeout');
}, 2000);

// clearInterval(id);

// clearTimeout(id2);

log('end');
