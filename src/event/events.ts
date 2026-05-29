import { type CalendarEventData, eventsEqual } from './event.js';

export type Events = CalendarEventData[];

function has(events: Events, target: CalendarEventData): boolean {
  return events.some(e => eventsEqual(e, target));
}

export interface DiffResult {
  toBeAdded: Events;
  toBeDeleted: Events;
}

export function diff(existing: Events, fetched: Events): DiffResult {
  const toBeAdded: Events = [];
  for (const e of fetched) {
    if (!has(existing, e) && !has(toBeAdded, e)) {
      toBeAdded.push(e);
    }
  }

  const toBeDeleted: Events = [];
  for (const e of existing) {
    if (!has(fetched, e) && !has(toBeDeleted, e)) {
      toBeDeleted.push(e);
    }
  }

  return { toBeAdded, toBeDeleted };
}
