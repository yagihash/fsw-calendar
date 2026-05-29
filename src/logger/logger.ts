export type LogLevel = 'DEBUG' | 'INFO' | 'WARN' | 'ERROR';

export interface Logger {
  debug(message: string, fields?: Record<string, unknown>): void;
  info(message: string, fields?: Record<string, unknown>): void;
  warn(message: string, fields?: Record<string, unknown>): void;
  error(message: string, fields?: Record<string, unknown>): void;
}

const levelPriority: Record<LogLevel, number> = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3,
};

function formatFields(fields?: Record<string, unknown>): string {
  if (!fields || Object.keys(fields).length === 0) return '';
  return ' ' + JSON.stringify(fields);
}

export function newLogger(level: LogLevel = 'INFO'): Logger {
  const minPriority = levelPriority[level];

  return {
    debug(message, fields) {
      if (minPriority <= levelPriority['DEBUG']) {
        process.stdout.write(`::debug::${message}${formatFields(fields)}\n`);
      }
    },
    info(message, fields) {
      if (minPriority <= levelPriority['INFO']) {
        console.log(`${message}${formatFields(fields)}`);
      }
    },
    warn(message, fields) {
      if (minPriority <= levelPriority['WARN']) {
        process.stdout.write(`::warning::${message}${formatFields(fields)}\n`);
      }
    },
    error(message, fields) {
      if (minPriority <= levelPriority['ERROR']) {
        process.stdout.write(`::error::${message}${formatFields(fields)}\n`);
      }
    },
  };
}
