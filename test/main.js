// const event = require('./test/eventEmitter.js');
// const timer = require('./test/timers.js');
// const fs = require('./test/fs.js');
// const processs = require('./test/process_envs.js');
// const queueMicroTask = require('./test/queueMicroTask.js');
// const tick = require('./test/nextTick.js');

// queueMicroTask;
// tick;
// event;
// timer;
// fs;
// processs;

console.log('Start');

const net = globalThis.net;

// Create server
const server = net.createServer((socket) => {
  console.log('[Server] Client connected');

  socket.on('data', (msg) => {
    console.log('[Server] Received:', msg);
    socket.write('Echo: ' + msg);
  });

  socket.on('end', () => {
    console.log('[Server] Client ended connection');
  });

  socket.on('close', () => {
    console.log('[Server] Socket closed');
  });
});

server.on('error', (err) => {
  console.log('[Server] Error:', err);
});

server.listen(5000, '127.0.0.1');

// Create client
const client = net.createConnection({ port: 5000, host: '127.0.0.1' }, () => {
  console.log('[Client] Connected to server');
  client.write('Hello Server');
});

client.on('data', (msg) => {
  console.log('[Client] Got:', msg);
  client.end(); // close after receiving echo
});

client.on('end', () => {
  console.log('[Client] Connection ended by server');
});

client.on('close', () => {
  console.log('[Client] Closed');
});
