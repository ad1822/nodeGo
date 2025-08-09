const startTime = Date.now();
// const fs = require('fs');

module.exports = function log(msg) {
  const elapsed = ((Date.now() - startTime) / 1000).toFixed(2);
  console.warn(`[+${elapsed}s] ${msg}`);
};
