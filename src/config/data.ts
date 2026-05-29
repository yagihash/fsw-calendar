import { type Course, parseCourse } from '../fetcher/course.js';
import { type Class, parseClass } from '../fetcher/class.js';

export interface Data {
  calendarId: string;
  course: Course;
  class: Class;
}

export function parseData(raw: string): Data {
  const obj = JSON.parse(raw) as Record<string, string>;
  return {
    calendarId: obj['calendar_id'] ?? '',
    course: parseCourse(obj['course'] ?? ''),
    class: parseClass(obj['class'] ?? ''),
  };
}
