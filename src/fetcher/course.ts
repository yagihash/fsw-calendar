export type Course = 'rc' | 'ss' | '';

export const RC: Course = 'rc';
export const SS: Course = 'ss';
export const Unknown: Course = '';

export function parseCourse(s: string): Course {
  switch (s) {
    case 'rc': return RC;
    case 'ss': return SS;
    default: return Unknown;
  }
}
