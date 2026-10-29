import { describe, it, expect } from 'vitest';
import { parseTarget, serializePayload } from './payload.js';
import { type Class } from '../fetcher/class.js';
import { type DocEvent } from '../event/event.js';

const e1: DocEvent = { date: '2026-10-29', start: '08:40', end: '09:10', title: 'R-4 A' };
const e2: DocEvent = { date: '2026-10-31', start: '09:00', end: '16:00', title: 'Sコース 4輪' };

describe('serializePayload / parseTarget', () => {
  const raw = serializePayload(2026, 10, new Map<Class, DocEvent[]>([['r-4', [e1]], ['s-course-4', [e2]]]));

  it('hands over the events of each class with the month', () => {
    expect(parseTarget('r-4', raw)).toEqual({ year: 2026, month: 10, events: [e1] });
    expect(parseTarget('s-course-4', raw)).toEqual({ year: 2026, month: 10, events: [e2] });
  });

  it('fills the classes which have no events with an empty array', () => {
    expect(parseTarget('n-4', raw)).toEqual({ year: 2026, month: 10, events: [] });
  });

  it('is a single line, which GITHUB_OUTPUT can take as it is', () => {
    expect(raw).not.toContain('\n');
  });
});

describe('parseTarget', () => {
  it('throws on invalid JSON', () => {
    expect(() => parseTarget('r-4', '{')).toThrow('failed to parse the events as JSON');
  });

  it('throws when the events of the class are missing', () => {
    expect(() => parseTarget('r-4', JSON.stringify({ year: 2026, month: 10, events: {} })))
      .toThrow('the events of r-4 are not found in the payload');
  });

  it('throws when the month is missing', () => {
    expect(() => parseTarget('r-4', JSON.stringify({ events: { 'r-4': [] } })))
      .toThrow('the events of r-4 are not found in the payload');
  });

  it('throws on an empty payload', () => {
    expect(() => parseTarget('r-4', '')).toThrow('failed to parse the events as JSON');
  });
});
