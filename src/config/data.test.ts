import { describe, it, expect } from 'vitest';
import { parseData } from './data.js';
import { R4 } from '../fetcher/class.js';

describe('parseData', () => {
  it('parses valid JSON correctly', () => {
    const input = JSON.stringify({
      calendar_id: 'https://example.com/test',
      class: 'r-4',
    });

    expect(parseData(input)).toEqual({
      calendarId: 'https://example.com/test',
      class: R4,
    });
  });

  it('sets unknown class for unrecognized value', () => {
    const input = JSON.stringify({ calendar_id: '', class: 'unknown' });
    expect(parseData(input).class).toBe('');
  });
});
