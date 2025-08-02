const startTime = Date.now();

function log(msg) {
  const elapsed = ((Date.now() - startTime) / 1000).toFixed(2);
  console.log(`[+${elapsed}s] ${msg}`);
}

log('Start');

let i = 0;
setInterval(() => {
  i++;

  log(`Interval ${i}`);
}, 1000);

// setTimeout(() => {
//   log('After 10 second');
// }, 10000);

// setTimeout(() => {
//   log('After 20 second');
// }, 20000);

setTimeout(() => {
  log('After 3 second');
}, 3000);

// setTimeout(() => {
//   log('After 2 second');
// }, 2000);

// setTimeout(() => {
//   log('After 1 second');
// }, 1000);

log('End');
