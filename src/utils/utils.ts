export function nextMonth(y: number, m: number): [number, number] {
  if (m === 12) {
    return [y + 1, 1];
  }
  return [y, m + 1];
}
