let id2 = setTimeout(() => {
  log('Timeout');
}, 2000);

setInterval(() => {
  if (i == 5) return;
  i++;
  console.warn(`Interval ${i}`);
}, 2000);

setImmediate(() => console.warn('immediate'));

setInterval(() => {
  log('Interval');
}, 1000);
