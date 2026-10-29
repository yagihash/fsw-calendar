import { describe, it, expect, vi, afterEach } from 'vitest';
import { fetchDocEvents } from './fetcher.js';
import { classes } from './class.js';

const pageHtml = '<script src="/s/sfsites/l/%7B%22fwuid%22%3A%22FWUID%22%2C%22loaded%22%3A%7B%22APPLICATION%40markup%3A%2F%2Fsiteforce%3AcommunityApp%22%3A%22LOADED%22%7D%7D/inline.js"></script>';

function slot(category: string, date: string, start: number, end: number, symbol?: string, canceled = false) {
  return {
    TMC_DX_DriveDate__c: date,
    TMC_DX_StartTime__c: start,
    TMC_DX_EndTime__c: end,
    TMC_DX_DriveSlotNumber__c: symbol,
    TMC_DX_IsReservationCanceled__c: canceled,
    TMC_DX_DriveCategory__r: { Name: category },
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

// posts are answered in order, one for each course
function stubFetch(...posts: Response[]) {
  const queue = [...posts];
  const fn = vi.fn((url: string) =>
    Promise.resolve(url.includes('/aura') ? queue.shift()! : response(200, pageHtml)),
  );
  vi.stubGlobal('fetch', fn);
  return fn;
}

function postParams(fn: ReturnType<typeof stubFetch>, call: number) {
  const [, init] = fn.mock.calls[call] as unknown as [string, RequestInit];
  const body = init.body as URLSearchParams;
  return {
    params: JSON.parse(body.get('message')!).actions[0].params.params,
    context: JSON.parse(body.get('aura.context')!),
  };
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('fetchDocEvents', () => {
  it('splits the slots of the classes in the same course', async () => {
    stubFetch(response(200, auraBody([
      slot('R-4', '2026-10-29', 33600000, 35400000, 'A'),
      slot('S-4', '2026-10-29', 31200000, 33000000, 'A'),
      slot('R-4', '2026-10-30', 33600000, 35400000, 'B'),
    ])));

    const got = await fetchDocEvents('example.com', ['r-4', 's-4', 'n-4'], 2026, 10);
    expect(got.get('r-4')).toEqual([
      { date: '2026-10-29', start: '09:20', end: '09:50', title: 'R-4 A' },
      { date: '2026-10-30', start: '09:20', end: '09:50', title: 'R-4 B' },
    ]);
    expect(got.get('s-4')).toEqual([
      { date: '2026-10-29', start: '08:40', end: '09:10', title: 'S-4 A' },
    ]);
    expect(got.get('n-4')).toEqual([]);
  });

  it('requests the page once and each course once for all the classes', async () => {
    const fn = stubFetch(response(200, auraBody([])), response(200, auraBody([])));

    await fetchDocEvents('example.com', classes, 2026, 3);

    expect(fn).toHaveBeenCalledTimes(3);
    expect(fn.mock.calls[0]![0]).toBe('https://example.com/s/racing-course-reservation?language=ja');
    expect(fn.mock.calls[1]![0]).toBe('https://example.com/s/sfsites/aura?r=1&aura.ApexAction.execute=1');

    const racing = postParams(fn, 1);
    expect(racing.params.driveCourse).toBe('レーシングコース');
    expect(racing.params.driveDate).toBe('2026-03-01');
    expect(JSON.parse(racing.params.driveCategoryList)).toEqual(['R-4', 'S-4', 'N-4', 'L-4']);
    expect(racing.context.fwuid).toBe('FWUID');
    expect(racing.context.loaded['APPLICATION@markup://siteforce:communityApp']).toBe('LOADED');

    const short = postParams(fn, 2);
    expect(short.params.driveCourse).toBe('ショートコース');
    expect(JSON.parse(short.params.driveCategoryList)).toEqual(['Ｓコース　４輪']);
  });

  it('omits the slot symbol when the site does not provide it', async () => {
    stubFetch(response(200, auraBody([slot('Ｓコース　４輪', '2026-10-31', 32400000, 57600000)])));

    const got = await fetchDocEvents('example.com', ['s-course-4'], 2026, 10);
    expect(got.get('s-course-4')).toEqual([
      { date: '2026-10-31', start: '09:00', end: '16:00', title: 'Sコース 4輪' },
    ]);
  });

  it('skips canceled slots and slots of unknown categories', async () => {
    stubFetch(response(200, auraBody([
      slot('R-4', '2026-10-29', 31200000, 33000000, 'A', true),
      slot('R-4', '2026-10-29', 33600000, 35400000, 'B'),
      slot('KART', '2026-10-29', 33600000, 35400000, 'A'),
    ])));

    const got = await fetchDocEvents('example.com', ['r-4'], 2026, 10);
    expect(got.get('r-4')!.map(e => e.title)).toEqual(['R-4 B']);
  });

  it('returns empty arrays when no slots are published', async () => {
    stubFetch(response(200, auraBody([])));

    const got = await fetchDocEvents('example.com', ['r-4'], 2026, 10);
    expect(got.get('r-4')).toEqual([]);
  });

  it('throws when the page does not have the aura context', async () => {
    vi.stubGlobal('fetch', () => Promise.resolve(response(200, '<html></html>')));
    await expect(fetchDocEvents('example.com', ['r-4'], 2026, 10))
      .rejects.toThrow('failed to find the aura context');
  });

  it('throws on non-200 from the page', async () => {
    vi.stubGlobal('fetch', () => Promise.resolve(response(503, '')));
    await expect(fetchDocEvents('example.com', ['r-4'], 2026, 10))
      .rejects.toThrow('got status code 503 on https://example.com/s/racing-course-reservation?language=ja');
  });

  it('throws on non-200 from the aura endpoint', async () => {
    stubFetch(response(500, ''));
    await expect(fetchDocEvents('example.com', ['r-4'], 2026, 10))
      .rejects.toThrow('got status code 500');
  });

  it('throws when the action fails', async () => {
    stubFetch(response(200, { actions: [{ state: 'ERROR', error: [{ message: 'no access' }] }] }));
    await expect(fetchDocEvents('example.com', ['r-4'], 2026, 10))
      .rejects.toThrow('failed to get drive slots (ERROR): no access');
  });
});
