import { register } from '../../src/index.js';

register().then(() => process.exit(0)).catch(err => {
  console.error(err);
  process.exit(1);
});
