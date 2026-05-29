import { google } from 'googleapis';
import { type CalendarEventData } from '../event/event.js';
import { type Events } from '../event/events.js';
import { nextMonth } from '../utils/utils.js';

function toJstDate(y: number, m: number, tz: string): string {
  return new Date(`${y}-${String(m).padStart(2, '0')}-01T00:00:00`).toLocaleString('sv-SE', {
    timeZone: tz,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).replace(' ', 'T') + '+09:00';
}

export class Calendar {
  private readonly calendarId: string;
  private readonly tz: string;
  private readonly service: ReturnType<typeof google.calendar>;

  constructor(calendarId: string, tz: string) {
    this.calendarId = calendarId;
    this.tz = tz;
    const auth = new google.auth.GoogleAuth({
      scopes: ['https://www.googleapis.com/auth/calendar'],
    });
    this.service = google.calendar({ version: 'v3', auth });
  }

  async getEvents(y: number, m: number, length: number): Promise<Events> {
    const timeMin = toJstDate(y, m, this.tz);
    let [endY, endM] = [y, m];
    for (let i = 0; i < length; i++) {
      [endY, endM] = nextMonth(endY, endM);
    }
    const timeMax = toJstDate(endY, endM, this.tz);

    const res = await this.service.events.list({
      calendarId: this.calendarId,
      showDeleted: false,
      singleEvents: true,
      timeMin,
      timeMax,
    });

    return (res.data.items ?? []).map(item => ({
      id: item.id ?? undefined,
      summary: item.summary ?? '',
      start: { dateTime: item.start?.dateTime ?? '' },
      end: { dateTime: item.end?.dateTime ?? '' },
    }));
  }

  async insert(event: CalendarEventData): Promise<void> {
    await this.service.events.insert({
      calendarId: this.calendarId,
      requestBody: {
        summary: event.summary,
        start: event.start,
        end: event.end,
      },
    });
  }

  async delete(eventId: string): Promise<void> {
    await this.service.events.delete({
      calendarId: this.calendarId,
      eventId,
    });
  }
}
