import { load } from './config/config.js';
import { calendarIds } from './config/calendars.js';
import { newLogger, type Logger, type LogLevel } from './logger/logger.js';
import { fetchDocEvents } from './fetcher/fetcher.js';
import { classes, classInfo, type Class } from './fetcher/class.js';
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

export async function register(): Promise<void> {
  const c = load();
  const log = newLogger(c.logLevel as LogLevel);
  const notify = new Slack(c.webhook);

  log.debug('logger is ready');

  const now = new Date(new Date().toLocaleString('en-US', { timeZone: c.timezone }));
  const y = now.getFullYear();
  const m = now.getMonth() + 1;

  let docEvents: Map<Class, DocEvent[]>;
  try {
    docEvents = await fetchDocEvents(c.hostname, classes, y, m);
  } catch (err) {
    log.error('failed to fetch schedules', { error: String(err) });
    await notify.warn(`error: ${String(err)}`).catch(() => undefined);
    throw err;
  }
  log.debug('loaded schedules', { events: Object.fromEntries(docEvents) });

  // a failure of a calendar must not block the others
  const failed: Class[] = [];
  for (const cls of classes) {
    try {
      await syncCalendar(log, notify, cls, docEvents.get(cls) ?? [], y, m);
    } catch (err) {
      failed.push(cls);
      log.error('failed to sync calendar', { class: cls, error: String(err) });
      await notify.warn(`error (${classInfo(cls).label}): ${String(err)}`).catch(() => undefined);
    }
  }

  if (failed.length > 0) {
    throw new Error(`failed to sync ${failed.length} calendar(s): ${failed.join(', ')}`);
  }
}
