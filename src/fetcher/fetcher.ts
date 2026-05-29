import * as cheerio from 'cheerio';
import { type Course } from './course.js';
import { type Class } from './class.js';
import { type DocEvent } from '../event/event.js';
import { nextMonth } from '../utils/utils.js';

function buildUrl(hostname: string, course: Course, cls: Class, y: number, m: number): string {
  return `https://${hostname}/driving/sports/${course}/${cls}/${y}/${String(m).padStart(2, '0')}.html`;
}

export async function fetchDocEvents(
  hostname: string,
  course: Course,
  cls: Class,
  y: number,
  m: number,
  length: number,
): Promise<DocEvent[]> {
  const results: DocEvent[] = [];
  let curY = y;
  let curM = m;

  for (let i = 0; i < length; i++) {
    const url = buildUrl(hostname, course, cls, curY, curM);
    const res = await fetch(url);

    if (i !== 0 && res.status === 404) {
      return results;
    }

    if (res.status !== 200) {
      throw new Error(`got status code ${res.status} on ${url}`);
    }

    const html = await res.text();
    const $ = cheerio.load(html);

    const typeElems = $('#table-calendar > tbody > tr.row-rc > td.type > div > p');
    const timeElems = $('#table-calendar > tbody > tr.row-rc > td.time > div > p');

    typeElems.each((idx, elem) => {
      const row = $(elem).parent().parent().parent();
      const date = row.attr('data-date') ?? '';
      const timeText = $(timeElems[idx]).text();
      const parts = timeText.split('~');
      if (parts.length < 2) return;
      const title = $(elem).text();
      results.push({ date, start: parts[0]!.trim(), end: parts[1]!.trim(), title });
    });

    [curY, curM] = nextMonth(curY, curM);
  }

  return results;
}
