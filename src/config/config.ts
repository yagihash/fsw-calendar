export interface Config {
  timezone: string;
  recurrence: number;
  logLevel: string;
  hostname: string;
  webhook: string;
}

export function load(): Config {
  return {
    timezone: process.env['TIMEZONE'] ?? 'Asia/Tokyo',
    recurrence: parseInt(process.env['RECURRENCE'] ?? '2', 10),
    logLevel: process.env['LOG_LEVEL'] ?? 'INFO',
    hostname: process.env['HOSTNAME'] ?? 'www.fsw.tv',
    webhook: process.env['SLACK_WEBHOOK'] ?? '',
  };
}
