export interface Notifier {
  info(message: string): Promise<void>;
  warn(message: string): Promise<void>;
}
