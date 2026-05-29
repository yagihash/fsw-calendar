import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { load } from './config.js';

describe('load', () => {
  let savedEnv: NodeJS.ProcessEnv;

  beforeEach(() => {
    savedEnv = { ...process.env };
    delete process.env['TIMEZONE'];
    delete process.env['RECURRENCE'];
    delete process.env['LOG_LEVEL'];
    delete process.env['HOSTNAME'];
    delete process.env['SLACK_WEBHOOK'];
  });

  afterEach(() => {
    process.env = savedEnv;
  });

  it('uses specified values', () => {
    process.env['TIMEZONE'] = 'TEST_TIMEZONE';
    process.env['RECURRENCE'] = '10';
    process.env['LOG_LEVEL'] = 'DEBUG';
    process.env['HOSTNAME'] = 'example.com';

    expect(load()).toEqual({
      timezone: 'TEST_TIMEZONE',
      recurrence: 10,
      logLevel: 'DEBUG',
      hostname: 'example.com',
      webhook: '',
    });
  });

  it('uses default timezone when not set', () => {
    expect(load().timezone).toBe('Asia/Tokyo');
  });

  it('uses default recurrence when not set', () => {
    expect(load().recurrence).toBe(2);
  });

  it('uses default log level when not set', () => {
    expect(load().logLevel).toBe('INFO');
  });

  it('uses default hostname when not set', () => {
    expect(load().hostname).toBe('www.fsw.tv');
  });
});
