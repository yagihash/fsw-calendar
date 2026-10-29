import { type Class, classes } from '../fetcher/class.js';
import { type DocEvent } from '../event/event.js';

// What the fetch job hands over to the register jobs.
export interface Payload {
  // the month which the events are fetched from (the first day of it)
  year: number;
  month: number;
  events: Record<Class, DocEvent[]>;
}

export function serializePayload(year: number, month: number, events: Map<Class, DocEvent[]>): string {
  const payload: Payload = {
    year,
    month,
    events: Object.fromEntries(classes.map(cls => [cls, events.get(cls) ?? []])) as Record<Class, DocEvent[]>,
  };
  return JSON.stringify(payload);
}

export interface Target {
  year: number;
  month: number;
  events: DocEvent[];
}

export function parseTarget(cls: Class, raw: string): Target {
  let payload: Partial<Payload>;
  try {
    payload = JSON.parse(raw) as Partial<Payload>;
  } catch {
    throw new Error('failed to parse the events as JSON');
  }

  const events = payload.events?.[cls];
  if (!Number.isInteger(payload.year) || !Number.isInteger(payload.month) || !Array.isArray(events)) {
    throw new Error(`the events of ${cls} are not found in the payload`);
  }

  return { year: payload.year!, month: payload.month!, events };
}
