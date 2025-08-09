const event = require('./test/eventEmitter.js');
const timer = require('./test/timers.js');
const fs = require('./test/fs.js');
const processs = require('./test/process_envs.js');
const queueMicroTask = require('./test/queueMicroTask.js');
const tick = require('./test/nextTick.js');

queueMicroTask;
tick;
event;
timer;
fs;
processs;
