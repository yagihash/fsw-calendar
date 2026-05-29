import { describe, it, expect } from 'vitest';
import { diff, type Events } from './events.js';
import { type CalendarEventData } from './event.js';

const t0 = new Date('2025-01-01T09:00:00+09:00');
const mins = (n: number) => new Date(t0.getTime() + n * 60000).toISOString();

const A: CalendarEventData = { summary: 'A', start: { dateTime: t0.toISOString() }, end: { dateTime: mins(25) } };
const B: CalendarEventData = { summary: 'B', start: { dateTime: mins(60) }, end: { dateTime: mins(85) } };
const C: CalendarEventData = { summary: 'C', start: { dateTime: mins(120) }, end: { dateTime: mins(145) } };
const D: CalendarEventData = { summary: 'D', start: { dateTime: mins(180) }, end: { dateTime: mins(205) } };

describe('diff', () => {
  it('identifies events to add and delete', () => {
    const existing: Events = [A, B, C];
    const fetched: Events = [A, B, D];

    const { toBeAdded, toBeDeleted } = diff(existing, fetched);

    expect(toBeAdded).toEqual([D]);
    expect(toBeDeleted).toEqual([C]);
  });

  it('returns empty arrays when events are identical', () => {
    const { toBeAdded, toBeDeleted } = diff([A, B], [A, B]);
    expect(toBeAdded).toHaveLength(0);
    expect(toBeDeleted).toHaveLength(0);
  });

  it('deduplicates within toBeAdded', () => {
    const { toBeAdded } = diff([], [A, A]);
    expect(toBeAdded).toHaveLength(1);
  });

  it('deduplicates within toBeDeleted', () => {
    const { toBeDeleted } = diff([A, A], []);
    expect(toBeDeleted).toHaveLength(1);
  });

  it('marks all fetched as toBeAdded when existing is empty', () => {
    const { toBeAdded, toBeDeleted } = diff([], [A, B]);
    expect(toBeAdded).toEqual([A, B]);
    expect(toBeDeleted).toHaveLength(0);
  });

  it('marks all existing as toBeDeleted when fetched is empty', () => {
    const { toBeAdded, toBeDeleted } = diff([A, B], []);
    expect(toBeAdded).toHaveLength(0);
    expect(toBeDeleted).toEqual([A, B]);
  });
});
