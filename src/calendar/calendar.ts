import { google } from 'googleapis';
import { type CalendarEventData } from '../event/event.js';
import { type Events } from '../event/events.js';

// The first moment of the month in JST. Events before it are never touched.
function monthStart(y: number, m: number): string {
  return `${y}-${String(m).padStart(2, '0')}-01T00:00:00+09:00`;
}

export class Calendar {
  private readonly calendarId: string;
  private readonly service: ReturnType<typeof google.calendar>;

  constructor(calendarId: string) {
    this.calendarId = calendarId;
    const auth = new google.auth.GoogleAuth({
      scopes: ['https://www.googleapis.com/auth/calendar'],
    });
    this.service = google.calendar({ version: 'v3', auth });
  }

  // Lists every event on and after the first day of the month, which is the range the fetcher covers.
  async getEvents(y: number, m: number): Promise<Events> {
    const res = await this.service.events.list({
      calendarId: this.calendarId,
      showDeleted: false,
      singleEvents: true,
      maxResults: 2500,
      timeMin: monthStart(y, m),
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
