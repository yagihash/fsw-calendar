import { load } from './config/config.js';
import { type Data } from './config/data.js';
import { newLogger, type LogLevel } from './logger/logger.js';
import { fetchDocEvents } from './fetcher/fetcher.js';
import { newFromDocEvent } from './event/event.js';
import { diff } from './event/events.js';
import { Calendar } from './calendar/calendar.js';
import { Slack } from './notify/slack/slack.js';

export async function register(data: Data): Promise<void> {
  const c = load();
  const log = newLogger(c.logLevel as LogLevel);
  const notify = new Slack(c.webhook);

  try {
    log.debug('logger is ready');

    const now = new Date(new Date().toLocaleString('en-US', { timeZone: c.timezone }));
    const y = now.getFullYear();
    const m = now.getMonth() + 1;

    const docEvents = await fetchDocEvents(c.hostname, data.course, data.class, y, m, c.recurrence);
    log.debug('loaded schedules', { events: docEvents });

    const fetchedEvents = docEvents.map(newFromDocEvent);

    const calendar = new Calendar(data.calendarId, c.timezone);
    const existingEvents = await calendar.getEvents(y, m, c.recurrence);

    const { toBeAdded, toBeDeleted } = diff(existingEvents, fetchedEvents);

    if (toBeAdded.length === 0 && toBeDeleted.length === 0) {
      log.debug('no update', { existing: existingEvents, fetched: fetchedEvents });
      return;
    }

    log.info('need updates', { to_be_added: toBeAdded, to_be_deleted: toBeDeleted });

    for (const e of toBeAdded) {
      try {
        await calendar.insert(e);
        log.debug('added new event', { event: e });
      } catch (err) {
        log.error('failed to insert event', { error: String(err), event: e });
      }
    }

    for (const e of toBeDeleted) {
      if (!e.id) continue;
      try {
        await calendar.delete(e.id);
        log.debug('deleted stale event', { event: e });
      } catch (err) {
        log.error('failed to delete event', { error: String(err), event: e });
      }
    }

    await notify.info(`updated calendar (${data.class.toUpperCase()})`);
  } catch (err) {
    log.error('unhandled error', { error: String(err) });
    await notify.warn(`error: ${String(err)}`).catch(() => undefined);
    throw err;
  }
}
