import { describe, it, expect } from 'vitest';
import { newEvent, eventsEqual } from './event.js';

describe('newEvent', () => {
  it('formats datetime with zero-padded hour', () => {
    const e = newEvent('2022-09-20', '07:00', '07:25', 'TEST_TITLE');
    expect(e.summary).toBe('TEST_TITLE');
    expect(e.start.dateTime).toBe('2022-09-20T07:00:00+09:00');
    expect(e.end.dateTime).toBe('2022-09-20T07:25:00+09:00');
  });

  it('zero-pads single-digit hour from FSW site', () => {
    const e = newEvent('2026-05-22', '9:30', '10:00', 'S-4 A');
    expect(e.start.dateTime).toBe('2026-05-22T09:30:00+09:00');
    expect(e.end.dateTime).toBe('2026-05-22T10:00:00+09:00');
  });
});

describe('eventsEqual', () => {
  const base = {
    summary: 'TEST_TITLE',
    start: { dateTime: '2022-09-20T07:00:00+09:00' },
    end: { dateTime: '2022-09-20T07:25:00+09:00' },
  };

  it('returns true for identical events', () => {
    expect(eventsEqual(base, { ...base })).toBe(true);
  });

  it('returns false when summary differs', () => {
    expect(eventsEqual(base, { ...base, summary: 'OTHER' })).toBe(false);
  });

  it('returns false when start differs', () => {
    const other = { ...base, start: { dateTime: '2022-09-19T7:00:00+09:00' } };
    expect(eventsEqual(base, other)).toBe(false);
  });

  it('returns false when end differs', () => {
    const other = { ...base, end: { dateTime: '2022-09-21T7:25:00+09:00' } };
    expect(eventsEqual(base, other)).toBe(false);
  });

  it('returns false for invalid (unparseable) datetime', () => {
    const broken = { ...base, start: { dateTime: '---' } };
    expect(eventsEqual(base, broken)).toBe(false);
  });
});
