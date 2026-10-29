import { appendFileSync } from 'node:fs';
import { fetchSchedules } from '../../src/index.js';

// Hands the schedules over to the register jobs through the output of the GitHub Actions job.
fetchSchedules().then(events => {
  const output = process.env['GITHUB_OUTPUT'];
  if (output) {
    // JSON.stringify never emits a raw newline, so the single-line form is safe
    appendFileSync(output, `events=${events}\n`);
  } else {
    console.log(events);
  }
  process.exit(0);
}).catch(err => {
  console.error(err);
  process.exit(1);
});
