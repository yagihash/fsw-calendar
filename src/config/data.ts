import { type Class, parseClass } from '../fetcher/class.js';

export interface Data {
  calendarId: string;
  class: Class;
}

export function parseData(raw: string): Data {
  const obj = JSON.parse(raw) as Record<string, string>;
  return {
    calendarId: obj['calendar_id'] ?? '',
    class: parseClass(obj['class'] ?? ''),
  };
}
