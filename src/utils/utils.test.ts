import { describe, it, expect } from 'vitest';
import { nextMonth } from './utils.js';

describe('nextMonth', () => {
  it.each<[number, number, number, number]>([
    [2000, 1, 2000, 2],
    [2000, 12, 2001, 1],
  ])('(%d, %d) → (%d, %d)', (y, m, ny, nm) => {
    const [gotY, gotM] = nextMonth(y, m);
    expect(gotY).toBe(ny);
    expect(gotM).toBe(nm);
  });
});
