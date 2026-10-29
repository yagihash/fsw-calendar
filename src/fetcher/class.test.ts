import { describe, it, expect } from 'vitest';
import { classes, classInfo, isClass } from './class.js';

describe('isClass', () => {
  it.each(classes)('accepts %s', c => {
    expect(isClass(c)).toBe(true);
  });

  it.each(['', 'ss-4', 'R-4', 'constructor'])('rejects %j', c => {
    expect(isClass(c)).toBe(false);
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

  it('has a distinct category for each class in a course', () => {
    const keys = classes.map(c => `${classInfo(c).course}/${classInfo(c).category}`);
    expect(new Set(keys).size).toBe(classes.length);
  });
});
