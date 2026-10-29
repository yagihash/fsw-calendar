import { register } from '../../src/index.js';
import { parseData } from '../../src/config/data.js';

async function run(): Promise<number> {
  const calendarId = process.env['CALENDAR_ID'];
  const cls = process.env['CLASS'];

  if (!calendarId || !cls) {
    console.error('CALENDAR_ID and CLASS environment variables are required');
    return 1;
  }

  const raw = JSON.stringify({ calendar_id: calendarId, class: cls });
  const data = parseData(raw);

  await register(data);
  return 0;
}

run().then(code => process.exit(code)).catch(err => {
  console.error(err);
  process.exit(1);
});
