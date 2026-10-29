import { register } from '../../src/index.js';

async function run(): Promise<number> {
  const cls = process.env['CLASS'];
  const events = process.env['EVENTS'];

  if (!cls || !events) {
    console.error('CLASS and EVENTS environment variables are required');
    return 1;
  }

  await register(cls, events);
  return 0;
}

run().then(code => process.exit(code)).catch(err => {
  console.error(err);
  process.exit(1);
});
