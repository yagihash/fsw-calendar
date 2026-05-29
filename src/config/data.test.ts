import { describe, it, expect } from 'vitest';
import { parseData } from './data.js';
import { RC } from '../fetcher/course.js';
import { T4 } from '../fetcher/class.js';

describe('parseData', () => {
  it('parses valid JSON correctly', () => {
    const input = JSON.stringify({
      calendar_id: 'https://example.com/test',
      course: 'rc',
      class: 't-4',
    });

    expect(parseData(input)).toEqual({
      calendarId: 'https://example.com/test',
      course: RC,
      class: T4,
    });
  });

  it('sets unknown course for unrecognized value', () => {
    const input = JSON.stringify({ calendar_id: '', course: 'unknown', class: 'ss-4' });
    expect(parseData(input).course).toBe('');
  });

  it('sets unknown class for unrecognized value', () => {
    const input = JSON.stringify({ calendar_id: '', course: 'rc', class: 'unknown' });
    expect(parseData(input).class).toBe('');
  });
});
