import { describe, it, expect } from 'vitest';
import { parseClass, classInfo, R4, S4, N4, L4, SCourse4, Unknown } from './class.js';

describe('parseClass', () => {
  it.each([
    ['r-4', R4],
    ['s-4', S4],
    ['n-4', N4],
    ['l-4', L4],
    ['s-course-4', SCourse4],
    ['ss-4', Unknown],
    ['', Unknown],
  ])('%s → %s', (input, want) => {
    expect(parseClass(input)).toBe(want);
  });
});

describe('classInfo', () => {
  it('maps racing course classes', () => {
    expect(classInfo('n-4')).toEqual({ course: 'レーシングコース', category: 'N-4', label: 'N-4' });
  });

  it('maps the short course class with full-width category name', () => {
    expect(classInfo('s-course-4')).toEqual({
      course: 'ショートコース',
      category: 'Ｓコース　４輪',
      label: 'Sコース 4輪',
    });
  });
});
