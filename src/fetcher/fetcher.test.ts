import { describe, it, expect, vi, afterEach } from 'vitest';
import { fetchDocEvents } from './fetcher.js';
import { RC } from './course.js';
import { T4 } from './class.js';

function mockHtml(date: string, title: string, time: string): string {
  return `<table id="table-calendar"><tbody>
    <tr class="row-rc" data-date="${date}">
      <td class="type"><div><p>${title}</p></div></td>
      <td class="time"><div><p>${time}</p></div></td>
    </tr>
  </tbody></table>`;
}

function makeResponse(status: number, body: string): Response {
  return { status, text: () => Promise.resolve(body) } as unknown as Response;
}

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('fetchDocEvents', () => {
  it('fetches events across multiple months', async () => {
    vi.stubGlobal('fetch', (url: string) => {
      if ((url as string).includes('/2023/11.html')) {
        return Promise.resolve(makeResponse(200, mockHtml('2023-11-02', 'T-4 X', '15:20~15:50')));
      }
      return Promise.resolve(makeResponse(200, mockHtml('2023-12-01', 'T-4 X', '15:20~15:50')));
    });

    const got = await fetchDocEvents('example.com', RC, T4, 2023, 11, 2);
    expect(got).toEqual([
      { date: '2023-11-02', start: '15:20', end: '15:50', title: 'T-4 X' },
      { date: '2023-12-01', start: '15:20', end: '15:50', title: 'T-4 X' },
    ]);
  });

  it('throws on 404 for the first month', async () => {
    vi.stubGlobal('fetch', () => Promise.resolve(makeResponse(404, '')));

    await expect(fetchDocEvents('example.com', RC, T4, 2023, 11, 1))
      .rejects.toThrow('got status code 404 on https://example.com/driving/sports/rc/t-4/2023/11.html');
  });

  it('stops gracefully on 404 for subsequent months', async () => {
    vi.stubGlobal('fetch', (url: string) => {
      if ((url as string).includes('/2023/11.html')) {
        return Promise.resolve(makeResponse(200, mockHtml('2023-11-02', 'T-4 X', '15:20~15:50')));
      }
      return Promise.resolve(makeResponse(404, ''));
    });

    const got = await fetchDocEvents('example.com', RC, T4, 2023, 11, 2);
    expect(got).toEqual([
      { date: '2023-11-02', start: '15:20', end: '15:50', title: 'T-4 X' },
    ]);
  });

  it('returns empty array for page with no events', async () => {
    vi.stubGlobal('fetch', () => Promise.resolve(makeResponse(200, '<html></html>')));

    const got = await fetchDocEvents('example.com', RC, T4, 2023, 11, 1);
    expect(got).toEqual([]);
  });

  it('skips rows where time has no tilde separator', async () => {
    const html = `<table id="table-calendar"><tbody>
      <tr class="row-rc" data-date="2023-11-02">
        <td class="type"><div><p>T-4 X</p></div></td>
        <td class="time"><div><p>終了</p></div></td>
      </tr>
    </tbody></table>`;
    vi.stubGlobal('fetch', () => Promise.resolve(makeResponse(200, html)));

    const got = await fetchDocEvents('example.com', RC, T4, 2023, 11, 1);
    expect(got).toEqual([]);
  });
});
