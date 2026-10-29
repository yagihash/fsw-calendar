export interface Config {
  timezone: string;
  logLevel: string;
  hostname: string;
  webhook: string;
}

export function load(): Config {
  return {
    timezone: process.env['TIMEZONE'] ?? 'Asia/Tokyo',
    logLevel: process.env['LOG_LEVEL'] ?? 'INFO',
    hostname: process.env['HOSTNAME'] ?? 'membership.fujimotorsportsforest.jp',
    webhook: process.env['SLACK_WEBHOOK'] ?? '',
  };
}
