import { IncomingWebhook } from '@slack/webhook';
import { type Notifier } from '../notify.js';

export class Slack implements Notifier {
  private readonly webhook: IncomingWebhook;

  constructor(webhookUrl: string) {
    this.webhook = new IncomingWebhook(webhookUrl);
  }

  async info(message: string): Promise<void> {
    await this.webhook.send({
      attachments: [
        {
          color: 'good',
          fallback: 'fsw-calendar info notification',
          text: message,
          mrkdwn_in: ['text'],
        },
      ],
    });
  }

  async warn(message: string): Promise<void> {
    await this.webhook.send({
      attachments: [
        {
          color: 'warning',
          fallback: 'fsw-calendar warning notification',
          text: message,
          mrkdwn_in: ['text'],
        },
      ],
    });
  }
}
