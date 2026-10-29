export const classes = ['r-4', 's-4', 'n-4', 'l-4', 's-course-4'] as const;

export type Class = (typeof classes)[number];

export function isClass(s: string): s is Class {
  return (classes as readonly string[]).includes(s);
}

export interface ClassInfo {
  // driveCourse parameter of the membership site
  course: string;
  // driveCategory name on the membership site
  category: string;
  // event title (a slot symbol is appended when the site provides it)
  label: string;
}

const racing = 'レーシングコース';
const short = 'ショートコース';

export function classInfo(cls: Class): ClassInfo {
  switch (cls) {
    case 'r-4': return { course: racing, category: 'R-4', label: 'R-4' };
    case 's-4': return { course: racing, category: 'S-4', label: 'S-4' };
    case 'n-4': return { course: racing, category: 'N-4', label: 'N-4' };
    case 'l-4': return { course: racing, category: 'L-4', label: 'L-4' };
    case 's-course-4': return { course: short, category: 'Ｓコース　４輪', label: 'Sコース 4輪' };
  }
}
