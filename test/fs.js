const fs = require('fs');
const message = 'Hello boys';
fs.writeFileSync('./test/output.txt', message, 'utf-8');
const data = fs.readFileSync('./test/output.txt', 'utf8');

// const data = fs.readFileSync('./test/test.txt', 'utf-8');
console.log('Content:', data);
