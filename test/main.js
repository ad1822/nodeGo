// const log = require('./test/logger.js');
// const fs = require('fs');
// const message = 'Hello boys';
// fs.writeFileSync('./test/output.txt', message, 'utf-8');
// const data = fs.readFileSync('./test/output.txt', 'utf8');
// // console.warn(data);

// // log('start');

// i = 0;

// setImmediate(() => console.warn('immediate'));

// Promise.resolve().then(() => {
//   console.error('Promise resolved');
// });

// // setInterval(() => {
// //   if (i == 5) return;
// //   i++;
// //   console.warn(`Interval ${i}`);
// // }, 2000);

// queueMicrotask(() => {
//   log('microtask 1');
// });

// process.nextTick(() => {
//   console.warn('NEXT TICK 1');
// });

// let id = setInterval(() => {
//   log('Interval');
// }, 1000);

// process.nextTick(() => {
//   console.warn('NEXT TICK 2');
// });

// // let id2 = setTimeout(() => {
// //   log('Timeout');
// // }, 2000);

// // log('end');

// let pro = process.cwd();
// console.log(process.env);
// console.log(process.argv);
// console.log('argv:', process.exit);

const e = new EventEmitter();

function listener(msg) {
  console.log('Received:', msg);
}

e.on('greet', listener);
e.emit('greet', 'Ayush'); // → Received: Ayush
// e.off('greet', listener);
e.emit('greet', 'again?'); // → (no output)
