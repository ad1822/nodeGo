# Net Module

- The Node.js `net` module provides an asynchronous network API, enabling the creation of stream-based TCP or IPC (Inter-Process Communication) servers and clients
- Code module, Essential for buildling applications that require direct network communication at the TCP level

## Key Functionalities

### 1. Creating TCP servers

- `net.createServer()` method allows to create TCP server that listens for incoming connections on a specified port and optionally a host
- The server emits a `connection` event when a new client connects, providing a `Socket` object representing the connection

### 2. Creating TCP Clients

- `net.createConnection()` method is used to establish a connection to a TCP server. It returns a `Socket` object that can be used to send and receive data over the established connection

### 3. Socket Object

- The `Socket` object represents a network socket and is a `Duplex` stream, meaning it can be both read from and written to. It provides events like `data`, `end`, `error`, and methods for sending and receiving data, controlling the connection, and accessing connection properties

### 4. Server Object

- The `Server` object returned by `net.createServer()` represents the TCP server. It provides methods for starting and stopping the server `listen(), close()` and emits events like `listening`, `connection`, and `close`

```js
const net = require('net');

const server = net.createServer((socket) => {
  console.log('Client connected');

  socket.on('data', (data) => {
    console.log(`Received from client: ${data}`);
    socket.write(`Echo: ${data}`); // Send data back to the client
  });

  socket.on('end', () => {
    console.log('Client disconnected');
  });

  socket.on('error', (err) => {
    console.error(`Socket error: ${err.message}`);
  });
});

server.listen(3000, () => {
  console.log('TCP server listening on port 3000');
});
```

- The net module is foundational for building custom network protocols, chat applications, or any scenario where direct control over TCP communication is required in Node.js.
