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

// An empty result against a non-empty calendar is more likely an outage or a spec change
// than every slot being withdrawn, so stop before wiping the calendar.
export function ensureFetched(existing: Events, fetched: Events): void {
  if (fetched.length === 0 && existing.length > 0) {
    throw new Error(`fetched no events while the calendar has ${existing.length} events; aborting to avoid deleting them all`);
  }
}
