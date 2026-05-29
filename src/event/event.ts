export interface CalendarEventData {
  id?: string;
  summary: string;
  start: { dateTime: string };
  end: { dateTime: string };
}

export interface DocEvent {
  date: string;
  start: string;
  end: string;
  title: string;
}

function formatDateTime(date: string, time: string): string {
  return `${date}T${time}:00+09:00`;
}

export function newEvent(date: string, start: string, end: string, title: string): CalendarEventData {
  return {
    summary: title,
    start: { dateTime: formatDateTime(date, start) },
    end: { dateTime: formatDateTime(date, end) },
  };
}

export function newFromDocEvent(d: DocEvent): CalendarEventData {
  return newEvent(d.date, d.start, d.end, d.title);
}

export function eventsEqual(a: CalendarEventData, b: CalendarEventData): boolean {
  const as = new Date(a.start.dateTime).getTime();
  const ae = new Date(a.end.dateTime).getTime();
  const bs = new Date(b.start.dateTime).getTime();
  const be = new Date(b.end.dateTime).getTime();
  return a.summary === b.summary && as === bs && ae === be;
}
