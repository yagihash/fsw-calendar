import { describe, it, expect } from 'vitest';
import { calendarIds } from './calendars.js';
import { classes } from '../fetcher/class.js';

describe('calendarIds', () => {
  it('has a distinct calendar for each class', () => {
    const ids = classes.map(c => calendarIds[c]);
    expect(ids.every(id => id.endsWith('@group.calendar.google.com'))).toBe(true);
    expect(new Set(ids).size).toBe(classes.length);
  });
});
