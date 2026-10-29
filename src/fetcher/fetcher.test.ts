import { describe, it, expect, vi, afterEach } from 'vitest';
import { fetchDocEvents } from './fetcher.js';

const pageHtml = '<script src="/s/sfsites/l/%7B%22fwuid%22%3A%22FWUID%22%2C%22loaded%22%3A%7B%22APPLICATION%40markup%3A%2F%2Fsiteforce%3AcommunityApp%22%3A%22LOADED%22%7D%7D/inline.js"></script>';

function slot(date: string, start: number, end: number, symbol: string | undefined, canceled = false) {
  return {
    TMC_DX_DriveDate__c: date,
    TMC_DX_StartTime__c: start,
    TMC_DX_EndTime__c: end,
    TMC_DX_DriveSlotNumber__c: symbol,
    TMC_DX_IsReservationCanceled__c: canceled,
  };
}

function auraBody(slots: unknown[]) {
  return { actions: [{ state: 'SUCCESS', returnValue: { returnValue: slots } }] };
}

function response(status: number, body: string | object): Response {
  return {
    status,
    text: () => Promise.resolve(body as string),
    json: () => Promise.resolve(body),
  } as unknown as Response;
}

function stubFetch(post: Response) {
  const fn = vi.fn((url: string) =>
    Promise.resolve(url.includes('/aura') ? post : response(200, pageHtml)),
  );
  vi.stubGlobal('fetch', fn);
  return fn;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('fetchDocEvents', () => {
  it('converts slots to doc events', async () => {
    stubFetch(response(200, auraBody([
      slot('2026-10-29', 31200000, 33000000, 'A'),
      slot('2026-10-30', 33600000, 35400000, 'B'),
    ])));

    const got = await fetchDocEvents('example.com', 'r-4', 2026, 10);
    expect(got).toEqual([
      { date: '2026-10-29', start: '08:40', end: '09:10', title: 'R-4 A' },
      { date: '2026-10-30', start: '09:20', end: '09:50', title: 'R-4 B' },
    ]);
  });

  it('omits the slot symbol when the site does not provide it', async () => {
    stubFetch(response(200, auraBody([slot('2026-10-31', 32400000, 57600000, undefined)])));

    const got = await fetchDocEvents('example.com', 's-course-4', 2026, 10);
    expect(got).toEqual([
      { date: '2026-10-31', start: '09:00', end: '16:00', title: 'Sコース 4輪' },
    ]);
  });

  it('skips canceled slots', async () => {
    stubFetch(response(200, auraBody([
      slot('2026-10-29', 31200000, 33000000, 'A', true),
      slot('2026-10-29', 33600000, 35400000, 'B'),
    ])));

    const got = await fetchDocEvents('example.com', 'r-4', 2026, 10);
    expect(got.map(e => e.title)).toEqual(['R-4 B']);
  });

  it('requests slots from the first day of the month for the class', async () => {
    const fn = stubFetch(response(200, auraBody([])));

    await fetchDocEvents('example.com', 's-course-4', 2026, 3);

    const [url, init] = fn.mock.calls[1] as unknown as [string, RequestInit];
    expect(url).toBe('https://example.com/s/sfsites/aura?r=1&aura.ApexAction.execute=1');

    const body = init.body as URLSearchParams;
    const params = JSON.parse(body.get('message')!).actions[0].params.params;
    expect(params.driveCourse).toBe('ショートコース');
    expect(params.driveDate).toBe('2026-03-01');
    expect(JSON.parse(params.driveCategoryList)).toEqual(['Ｓコース　４輪']);

    const context = JSON.parse(body.get('aura.context')!);
    expect(context.fwuid).toBe('FWUID');
    expect(context.loaded['APPLICATION@markup://siteforce:communityApp']).toBe('LOADED');
  });

  it('returns empty array when no slots are published', async () => {
    stubFetch(response(200, auraBody([])));
    expect(await fetchDocEvents('example.com', 'r-4', 2026, 10)).toEqual([]);
  });

  it('throws when the page does not have the aura context', async () => {
    vi.stubGlobal('fetch', () => Promise.resolve(response(200, '<html></html>')));
    await expect(fetchDocEvents('example.com', 'r-4', 2026, 10))
      .rejects.toThrow('failed to find the aura context');
  });

  it('throws on non-200 from the page', async () => {
    vi.stubGlobal('fetch', () => Promise.resolve(response(503, '')));
    await expect(fetchDocEvents('example.com', 'r-4', 2026, 10))
      .rejects.toThrow('got status code 503 on https://example.com/s/racing-course-reservation?language=ja');
  });

  it('throws on non-200 from the aura endpoint', async () => {
    stubFetch(response(500, ''));
    await expect(fetchDocEvents('example.com', 'r-4', 2026, 10))
      .rejects.toThrow('got status code 500');
  });

  it('throws when the action fails', async () => {
    stubFetch(response(200, { actions: [{ state: 'ERROR', error: [{ message: 'no access' }] }] }));
    await expect(fetchDocEvents('example.com', 'r-4', 2026, 10))
      .rejects.toThrow('failed to get drive slots (ERROR): no access');
  });
});
