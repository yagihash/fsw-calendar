export type Class = 'ss-4' | 't-4' | 'ns-4' | 's-4' | '';

export const SS4: Class = 'ss-4';
export const T4: Class = 't-4';
export const NS4: Class = 'ns-4';
export const S4: Class = 's-4';
export const Unknown: Class = '';

export function parseClass(s: string): Class {
  switch (s) {
    case 'ss-4': return SS4;
    case 't-4': return T4;
    case 'ns-4': return NS4;
    case 's-4': return S4;
    default: return Unknown;
  }
}
