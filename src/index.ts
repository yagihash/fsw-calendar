import { load } from './config/config.js';
import { calendarIds } from './config/calendars.js';
import { parseTarget, serializePayload } from './config/payload.js';
import { newLogger, type Logger, type LogLevel } from './logger/logger.js';
import { fetchDocEvents } from './fetcher/fetcher.js';
import { classes, classInfo, isClass, type Class } from './fetcher/class.js';
import { type DocEvent, newFromDocEvent } from './event/event.js';
import { diff, ensureFetched } from './event/events.js';
import { Calendar } from './calendar/calendar.js';
import { Slack } from './notify/slack/slack.js';

async function syncCalendar(
  log: Logger,
  notify: Slack,
  cls: Class,
  docEvents: DocEvent[],
  y: number,
  m: number,
): Promise<void> {
  const fetchedEvents = docEvents.map(newFromDocEvent);

  const calendar = new Calendar(calendarIds[cls]);
  const existingEvents = await calendar.getEvents(y, m);
  ensureFetched(existingEvents, fetchedEvents);

  const { toBeAdded, toBeDeleted } = diff(existingEvents, fetchedEvents);

  if (toBeAdded.length === 0 && toBeDeleted.length === 0) {
    log.debug('no update', { class: cls, existing: existingEvents, fetched: fetchedEvents });
    return;
  }

  log.info('need updates', { class: cls, to_be_added: toBeAdded, to_be_deleted: toBeDeleted });

  for (const e of toBeAdded) {
    try {
      await calendar.insert(e);
      log.debug('added new event', { class: cls, event: e });
    } catch (err) {
      log.error('failed to insert event', { class: cls, error: String(err), event: e });
    }
  }

  for (const e of toBeDeleted) {
    if (!e.id) continue;
    try {
      await calendar.delete(e.id);
      log.debug('deleted stale event', { class: cls, event: e });
    } catch (err) {
      log.error('failed to delete event', { class: cls, error: String(err), event: e });
    }
  }

  await notify.info(`updated calendar (${classInfo(cls).label})`);
}

// The first day of the current month, which the schedules are fetched from.
function currentMonth(timezone: string): { year: number; month: number } {
  const now = new Date(new Date().toLocaleString('en-US', { timeZone: timezone }));
  return { year: now.getFullYear(), month: now.getMonth() + 1 };
}

// Fetches the schedules of every class at once and returns them serialized for the register jobs.
export async function fetchSchedules(): Promise<string> {
  const c = load();
  const log = newLogger(c.logLevel as LogLevel);
  const notify = new Slack(c.webhook);

  log.debug('logger is ready');

  const { year, month } = currentMonth(c.timezone);

  try {
    const docEvents = await fetchDocEvents(c.hostname, classes, year, month);
    log.debug('loaded schedules', { events: Object.fromEntries(docEvents) });
    return serializePayload(year, month, docEvents);
  } catch (err) {
    log.error('failed to fetch schedules', { error: String(err) });
    await notify.warn(`error: ${String(err)}`).catch(() => undefined);
    throw err;
  }
}

// Syncs the calendar of the class with the schedules fetched by fetchSchedules.
export async function register(cls: string, raw: string): Promise<void> {
  const c = load();
  const log = newLogger(c.logLevel as LogLevel);
  const notify = new Slack(c.webhook);

  try {
    log.debug('logger is ready');

    if (!isClass(cls)) {
      throw new Error(`unknown class: ${JSON.stringify(cls)}`);
    }

    const { year, month, events } = parseTarget(cls, raw);
    await syncCalendar(log, notify, cls, events, year, month);
  } catch (err) {
    log.error('failed to sync calendar', { class: cls, error: String(err) });
    await notify.warn(`error (${cls}): ${String(err)}`).catch(() => undefined);
    throw err;
  }
}
