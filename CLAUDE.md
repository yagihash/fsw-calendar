# fsw-calendar

富士スピードウェイ（FSW）の会員向け走行枠スケジュールを自動的にGoogleカレンダーに登録するTypeScript実装です。

## プロジェクト概要

富士モータースポーツフォレストの会員サイトからスケジュールを取得し、Googleカレンダーと差分同期します。GitHub Actionsで毎朝自動実行されます。

## ディレクトリ構成

- `src/index.ts` - メインエントリーポイント（`register`関数）
- `src/calendar/` - Google Calendar APIクライアント
- `src/config/` - 環境変数・設定読み込み、入力パース
- `src/event/` - イベントモデルと差分ロジック
- `src/fetcher/` - FSWサイトからのスケジュール取得
- `src/logger/` - GitHub Actions ワークフローコマンド準拠のロガー
- `src/notify/slack/` - Slackへのエラー通知
- `cmd/fetch/` - 全クラスのスケジュールを取得するCLIエントリーポイント（GitHub Actions の fetch ジョブから呼び出す）
- `cmd/run/` - 1クラス分のカレンダーを同期するCLIエントリーポイント（GitHub Actions の register ジョブから呼び出す）

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
2. fetch ジョブ: `cmd/fetch/main.ts` が会員サイトから全クラス（R-4 / S-4 / N-4 / L-4 / Sコース4輪）の当月1日以降の公開済みスケジュールをまとめて取得する（サイトへのリクエストは、ページ1回＋コースごとに1回）。結果は取得した年月と一緒に1行のJSONにして、ジョブの output（`events`）で渡す
3. register ジョブ: クラスごとの matrix（`fail-fast: false`）。`cmd/run/main.ts` が環境変数（`CLASS`, `EVENTS`）を読み込んで `register()` を呼び出す
4. 既存のGoogleカレンダーイベントと差分比較
5. 追加・削除が必要なイベントのみ更新
6. エラー発生時はSlackに通知。fetch が失敗した場合は register ジョブは走らない

fetch ジョブは Google の認証を持たず、register ジョブはサイトにアクセスしません。ジョブ間で渡すJSONの形式は `src/config/payload.ts` にあります。job output の上限は1MBで、現在は約12KBです。

## 対応カレンダー

- R-4
- S-4
- N-4
- L-4
- Sコース 4輪（ショートコース）

## 環境変数

`src/config/config.ts` の `Config` インターフェースで管理されています。クラスごとの Google カレンダー ID は、公開情報なので `src/config/calendars.ts` に置いています。

| 変数名 | デフォルト値 | 必須 | 説明 |
|---|---|---|---|
| `TIMEZONE` | `Asia/Tokyo` | | タイムゾーン |
| `LOG_LEVEL` | `INFO` | | ログレベル（DEBUG/INFO/WARN/ERROR） |
| `HOSTNAME` | `membership.fujimotorsportsforest.jp` | | 取得先の会員サイトのホスト名 |
| `SLACK_WEBHOOK` | — | ○ | Slack Incoming Webhook URL |

## スケジュール取得

会員サイトは Salesforce Experience Cloud（Aura）製の SPA で、HTML にスケジュールは含まれません。画面が内部で呼んでいる Apex を直接呼びます（ログイン不要、`src/fetcher/fetcher.ts`）。

1. `https://{HOSTNAME}/s/racing-course-reservation?language=ja` を取得し、HTML から `fwuid` と `loaded`（Aura のバージョン）を取り出す。リリースごとに変わるので毎回取得する
2. `POST https://{HOSTNAME}/s/sfsites/aura?r=1&aura.ApexAction.execute=1` で `TMC_DX_DriveSlotController.getDriveSlotList` を呼ぶ

主なパラメータ:

- `driveCourse`: `レーシングコース` または `ショートコース`。1回の呼び出しで指定できるコースは1つ（`src/fetcher/class.ts` の `classInfo()` でクラスから決まる）
- `driveCategoryList`: 取得するクラス名の JSON 配列。同じコースのクラスはまとめて指定し、レスポンスの `TMC_DX_DriveCategory__r.Name` でクラスに振り分ける。ショートコースは全角（`Ｓコース　４輪`、スペースも全角）
- `driveDate`: 取得開始日。当月1日を指定する
- `driveDateList`: `[]` にすると公開済みの全枠が返る。`maxDays` は `driveDate` からの日数上限

レスポンスの時刻（`TMC_DX_StartTime__c` など）は0時からのミリ秒です。キャンセル済み（`TMC_DX_IsReservationCanceled__c`）の枠は除外します。残数はイベント名に含めません（毎日差分が出るため）。

サイトのデータは最近のリプレース以降のものしかなく、過去の月は取得できません。

## コード設計の注意点

**イベント同一性の比較**
`eventsEqual()` はSummary・Start・End の3フィールドで同一性を判断します。Google Calendar のイベントIDでは比較しません（`src/event/event.ts`）。

**差分ロジック**
`diff(existing, fetched)` は集合演算で実装されています。`fetched` に存在して `existing` にないものが「追加すべきイベント」（`toBeAdded`）、`existing` に存在して `fetched` にないものが「削除すべきイベント」（`toBeDeleted`）になります（`src/event/events.ts`）。

**比較範囲**
取得側もカレンダー側も「当月1日 00:00 JST 以降すべて」を範囲とします。月は fetch ジョブが決めて JSON に含め、register ジョブはそれを使います（ジョブ間で月がずれて、過去分が重複登録されるのを防ぐため）。範囲が揃っていないと、範囲外のイベントが差分で消えたり重複したりします。前月以前のイベントは触りません（`src/calendar/calendar.ts`）。

**空取得の安全策**
取得結果が0件なのにカレンダーにイベントがある場合は、削除せずエラーで止めてSlackに通知します（`ensureFetched()`、`src/event/events.ts`）。

**時刻フォーマット**
`formatDateTime()` で `9:30` のような1桁の時刻もゼロ埋めして ISO 8601 準拠の形式（`09:30`）に統一してから Google Calendar API と比較する（`src/event/event.ts`）。

**Google Calendar 認証**
Application Default Credentials（ADC）を使用します。ローカルで動かす場合は `gcloud auth application-default login` が必要です。GitHub Actions では Workload Identity Federation で認証します。

**ロガー**
GitHub Actions のワークフローコマンドを使用します。`INFO` → `console.log`、`DEBUG` → `::debug::`、`WARN` → `::warning::`、`ERROR` → `::error::`。

## ローカル実行

```bash
EVENTS=$(pnpm exec tsx ./cmd/fetch/main.ts) \
  CLASS=n-4 \
  SLACK_WEBHOOK=https://hooks.slack.com/... \
  pnpm exec tsx ./cmd/run/main.ts
```

`cmd/fetch/main.ts` は `GITHUB_OUTPUT` がなければ JSON を標準出力に出します。`cmd/run/main.ts` は、そのJSONを `EVENTS`、クラスを `CLASS`（`r-4` / `s-4` / `n-4` / `l-4` / `s-course-4`）で渡すと、そのクラスの実カレンダーを更新します。

## テスト方針

- 全モジュールにユニットテストを記述する（Vitest）
- `src/fetcher/fetcher.test.ts` は `vi.stubGlobal('fetch', ...)` でHTTPクライアントをモックし、Aura のレスポンスを返す
- Google Calendar APIは結合テストなし（実APIを直接テストしない）

## Git ワークフロー

- `main` ブランチには直接 push しない
- 必ず feature ブランチを切って PR 経由でマージする

## 注意事項

- このカレンダーはFSW公式情報に基づく非公式な利用であり、可用性は保証されない
- スケジュール変更は最大24時間遅れて反映される
