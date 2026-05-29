# fsw-calendar

富士スピードウェイ（FSW）の会員向け走行枠スケジュールを自動的にGoogleカレンダーに登録するTypeScript実装です。

## プロジェクト概要

FSWの公式サイトからスケジュールを取得し、Googleカレンダーと差分同期します。GitHub Actionsで毎朝自動実行されます。

## ディレクトリ構成

- `src/index.ts` - メインエントリーポイント（`register`関数）
- `src/calendar/` - Google Calendar APIクライアント
- `src/config/` - 環境変数・設定読み込み、入力パース
- `src/event/` - イベントモデルと差分ロジック
- `src/fetcher/` - FSWサイトからのスケジュール取得
- `src/logger/` - GitHub Actions ワークフローコマンド準拠のロガー
- `src/notify/slack/` - Slackへのエラー通知
- `src/utils/` - 月インクリメントなどのヘルパー
- `cmd/run/` - CLIエントリーポイント（GitHub Actions から呼び出す）

## 開発コマンド

```bash
# Node.jsバージョンのインストール
mise install

# 依存関係のインストール
mise run setup

# テスト実行
mise run test

# 型チェック
mise run typecheck
```

## アーキテクチャ

1. GitHub Actions のcronスケジュール（毎日 17:00 UTC）で起動
2. matrix で SS-4 / T-4 / NS-4 / S-4 を並列実行
3. `cmd/run/main.ts` が環境変数（`CALENDAR_ID`, `COURSE`, `CLASS`）を読み込んで `register()` を呼び出す
4. FSWサイトから該当月のスケジュールをスクレイピング
5. 既存のGoogleカレンダーイベントと差分比較
6. 追加・削除が必要なイベントのみ更新
7. エラー発生時はSlackに通知

## 対応カレンダー

- SS-4
- T-4
- NS-4
- S-4

## 環境変数

`src/config/config.ts` の `Config` インターフェースで管理されています。

| 変数名 | デフォルト値 | 必須 | 説明 |
|---|---|---|---|
| `TIMEZONE` | `Asia/Tokyo` | | タイムゾーン |
| `RECURRENCE` | `2` | | 取得する月数（当月含む） |
| `LOG_LEVEL` | `INFO` | | ログレベル（DEBUG/INFO/WARN/ERROR） |
| `HOSTNAME` | `www.fsw.tv` | | スクレイピング先FSWのホスト名 |
| `SLACK_WEBHOOK` | — | ○ | Slack Incoming Webhook URL |
| `CALENDAR_ID` | — | ○ | Google カレンダー ID |
| `COURSE` | — | ○ | コース（`rc` / `ss`） |
| `CLASS` | — | ○ | クラス（`ss-4` / `t-4` / `ns-4` / `s-4`） |

## スクレイピング

FSWサイトから以下のURL形式でHTMLを取得し、cheerio でパースします。

```
https://{HOSTNAME}/driving/sports/{course}/{class}/{year}/{month:02d}.html
```

CSSセレクタ: `#table-calendar > tbody > tr.row-rc > td.type > div > p`

未来月のページが未公開（HTTP 404）の場合は正常として扱い、処理を終了します（`src/fetcher/fetcher.ts`参照）。

## コード設計の注意点

**イベント同一性の比較**
`eventsEqual()` はSummary・Start・End の3フィールドで同一性を判断します。Google Calendar のイベントIDでは比較しません（`src/event/event.ts`）。

**差分ロジック**
`diff(existing, fetched)` は集合演算で実装されています。`fetched` に存在して `existing` にないものが「追加すべきイベント」（`toBeAdded`）、`existing` に存在して `fetched` にないものが「削除すべきイベント」（`toBeDeleted`）になります（`src/event/events.ts`）。

**時刻フォーマット**
FSWサイトは `9:30` のように1桁の時刻を返すことがある。`formatDateTime()` でゼロ埋めして ISO 8601 準拠の形式（`09:30`）に統一してから Google Calendar API と比較する（`src/event/event.ts`）。

**Google Calendar 認証**
Application Default Credentials（ADC）を使用します。ローカルで動かす場合は `gcloud auth application-default login` が必要です。GitHub Actions では Workload Identity Federation で認証します。

**ロガー**
GitHub Actions のワークフローコマンドを使用します。`INFO` → `console.log`、`DEBUG` → `::debug::`、`WARN` → `::warning::`、`ERROR` → `::error::`。

## ローカル実行

```bash
CALENDAR_ID=xxx@group.calendar.google.com \
  COURSE=rc \
  CLASS=ns-4 \
  SLACK_WEBHOOK=https://hooks.slack.com/... \
  pnpm exec tsx ./cmd/run/main.ts
```

## テスト方針

- 全モジュールにユニットテストを記述する（Vitest）
- `src/fetcher/fetcher.test.ts` は `vi.stubGlobal('fetch', ...)` でHTTPクライアントをモックする
- Google Calendar APIは結合テストなし（実APIを直接テストしない）

## 注意事項

- このカレンダーはFSW公式情報に基づく非公式な利用であり、可用性は保証されない
- スケジュール変更は最大24時間遅れて反映される
