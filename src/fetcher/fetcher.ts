import { type Class, classInfo } from './class.js';
import { type DocEvent } from '../event/event.js';

const accountName = '富士モータースポーツフォレスト';

interface DriveSlot {
  TMC_DX_DriveDate__c: string;
  TMC_DX_StartTime__c: number;
  TMC_DX_EndTime__c: number;
  TMC_DX_DriveSlotNumber__c?: string;
  TMC_DX_IsReservationCanceled__c: boolean;
}

interface AuraResponse {
  actions?: {
    state: string;
    returnValue?: { returnValue?: DriveSlot[] };
    error?: { message?: string }[];
  }[];
}

// The membership site is a Salesforce Experience Cloud (Aura) app.
// fwuid and the app version change on every release, so pick them up from the page each time.
async function fetchAuraContext(pageUrl: string): Promise<Record<string, unknown>> {
  const res = await fetch(pageUrl);
  if (res.status !== 200) {
    throw new Error(`got status code ${res.status} on ${pageUrl}`);
  }
  const html = await res.text();

  const fwuid = /fwuid%22%3A%22([^%]+)%22/.exec(html)?.[1];
  const loaded = /APPLICATION%40markup%3A%2F%2Fsiteforce%3AcommunityApp%22%3A%22([^%]+)%22/.exec(html)?.[1];
  if (!fwuid || !loaded) {
    throw new Error(`failed to find the aura context on ${pageUrl}`);
  }

  return {
    mode: 'PROD',
    fwuid,
    app: 'siteforce:communityApp',
    loaded: { 'APPLICATION@markup://siteforce:communityApp': loaded },
    dn: [],
    globals: {},
    uad: true,
  };
}

// milliseconds since midnight -> HH:MM
function formatTime(ms: number): string {
  const minutes = Math.floor(ms / 60000);
  const h = String(Math.floor(minutes / 60)).padStart(2, '0');
  const m = String(minutes % 60).padStart(2, '0');
  return `${h}:${m}`;
}

// Fetches every published slot of the class on and after the first day of the given month.
export async function fetchDocEvents(
  hostname: string,
  cls: Exclude<Class, ''>,
  y: number,
  m: number,
): Promise<DocEvent[]> {
  const info = classInfo(cls);
  const pageUrl = `https://${hostname}/s/racing-course-reservation?language=ja`;
  const context = await fetchAuraContext(pageUrl);

  const message = {
    actions: [{
      id: '1;a',
      descriptor: 'aura://ApexActionController/ACTION$execute',
      callingDescriptor: 'UNKNOWN',
      params: {
        namespace: '',
        classname: 'TMC_DX_DriveSlotController',
        method: 'getDriveSlotList',
        params: {
          accountName,
          driveCourse: info.course,
          driveDate: `${y}-${String(m).padStart(2, '0')}-01`,
          driveCategoryList: JSON.stringify([info.category]),
          driveDateList: '[]',
          maxDays: 1000,
          limitNumber: 5000,
        },
        cacheable: false,
        isContinuation: false,
      },
    }],
  };

  const url = `https://${hostname}/s/sfsites/aura?r=1&aura.ApexAction.execute=1`;
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'content-type': 'application/x-www-form-urlencoded;charset=UTF-8' },
    body: new URLSearchParams({
      message: JSON.stringify(message),
      'aura.context': JSON.stringify(context),
      'aura.pageURI': '/s/racing-course-reservation?language=ja',
      'aura.token': 'null',
    }),
  });
  if (res.status !== 200) {
    throw new Error(`got status code ${res.status} on ${url}`);
  }

  const action = ((await res.json()) as AuraResponse).actions?.[0];
  if (action?.state !== 'SUCCESS') {
    const reason = action?.error?.map(e => e.message).join(', ') ?? 'no action in the response';
    throw new Error(`failed to get drive slots (${action?.state ?? 'UNKNOWN'}): ${reason}`);
  }

  return (action.returnValue?.returnValue ?? [])
    .filter(s => !s.TMC_DX_IsReservationCanceled__c)
    .map(s => ({
      date: s.TMC_DX_DriveDate__c,
      start: formatTime(s.TMC_DX_StartTime__c),
      end: formatTime(s.TMC_DX_EndTime__c),
      title: s.TMC_DX_DriveSlotNumber__c ? `${info.label} ${s.TMC_DX_DriveSlotNumber__c}` : info.label,
    }));
}
