const startTime = Date.now();
const fs = require('fs');

module.exports = function log(msg) {
  const elapsed = ((Date.now() - startTime) / 1000).toFixed(2);
  console.log(`[+${elapsed}s] ${msg}`);
};

const data = fs.readFileSync('./test/test.txt', 'utf-8');
console.log('Content:', data);
