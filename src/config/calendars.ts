import { type Class } from '../fetcher/class.js';

// Google Calendar IDs which each class is synced to. The calendars are public, so the IDs are not secrets.
export const calendarIds: Record<Class, string> = {
  'r-4': 'a7d264798dfbd87a96929ea4ae415b9c900a645b73403835a6b1621dbbd668cd@group.calendar.google.com',
  's-4': '4fca9ba5c00115af8aed9cbc340772ea5c5e64b3dfa975732f6e767862f405b1@group.calendar.google.com',
  'n-4': 'bdb6d8a80c9cdfd605319bc879d626d79fda0394cdeb231c25a7a4278f6f28ef@group.calendar.google.com',
  'l-4': '0e51d284e64c42f6067a27721b921e1dcb4759f04bbed5cf4fd4faff2e20071b@group.calendar.google.com',
  's-course-4': 'b78abe03924c2624a1c0f4271fcb6eb2a90a68f2029f239763ce06df1fb711fb@group.calendar.google.com',
};
