# fsw-calendar

富士スピードウェイ（FSW）の会員向け走行枠スケジュールを自動的にGoogleカレンダーに登録するGCP Cloud Functions実装です。

## プロジェクト概要

FSWの公式サイトからスケジュールを取得し、Googleカレンダーと差分同期します。Pub/Subメッセージで起動し、毎朝自動実行されます。

## ディレクトリ構成

- `function.go` - Cloud Functionsエントリーポイント（`Register`関数）
- `calendar/` - Google Calendar APIクライアント
- `config/` - 環境変数・設定読み込み
- `event/` - イベントモデルと差分ロジック
- `fetcher/` - FSWサイトからのスケジュール取得
- `fetcher/course/` - コース種別enum（RC/SS）
- `fetcher/class/` - クラス種別enum（SS4/T4/NS4/S4）
- `logger/` - zapベースのロガー初期化
- `notify/slack/` - Slackへのエラー通知
- `pages/` - テストカバレッジレポート（GitHub Pages）
- `cmd/reset/` - カレンダーイベント一括削除ツール（ローカル用）
- `utils/` - 月インクリメントなどのヘルパー

## 開発コマンド

```bash
# Goバージョンのインストール
mise install

# 依存関係のインストール
mise run setup

# テスト実行
mise run test

# カバレッジレポート生成・表示
mise run coverage
```

## アーキテクチャ

1. Cloud SchedulerがPub/SubトピックにメッセージをPublish
2. `Register`関数がPub/Subメッセージを受信
3. メッセージからコース・クラス・カレンダーIDを取得
4. FSWサイトから該当月のスケジュールをスクレイピング
5. 既存のGoogleカレンダーイベントと差分比較
6. 追加・削除が必要なイベントのみ更新
7. エラー発生時はSlackに通知

## Pub/Sub メッセージスキーマ

Cloud Schedulerから送られるメッセージのJSONフォーマット：

```json
{
  "calendar_id": "xxx@group.calendar.google.com",
  "course": "rc",
  "class": "ss-4"
}
```

- `course`: `"rc"` (レーシングコース) / `"ss"` (ショートサーキット)
- `class`: `"ss-4"` / `"t-4"` / `"ns-4"` / `"s-4"`（小文字ハイフン区切り）

## 対応カレンダー

- SS-4
- T-4
- NS-4
- S-4

## 環境変数

`config/config.go` の `Config` 構造体で管理されています。

| 変数名 | デフォルト値 | 必須 | 説明 |
|---|---|---|---|
| `TIMEZONE` | `Asia/Tokyo` | | タイムゾーン |
| `RECURRENCE` | `2` | | 取得する月数（当月含む） |
| `LOG_LEVEL` | `INFO` | | ログレベル（DEBUG/INFO/WARN/ERROR/FATAL） |
| `HOSTNAME` | `www.fsw.tv` | | スクレイピング先FSWのホスト名 |
| `SLACK_WEBHOOK` | — | ○ | Slack Incoming Webhook URL |
| `CALENDAR_ID` | — | ○（cmd/resetのみ） | Google カレンダー ID |

## スクレイピング

FSWサイトから以下のURL形式でHTMLを取得し、goquery でパースします。

```
https://{HOSTNAME}/driving/sports/{course}/{class}/{year}/{month:02d}.html
```

CSSセレクタ: `#table-calendar > tbody > tr.row-rc > td.type > div > p`

未来月のページが未公開（HTTP 404）の場合は正常として扱い、処理を終了します（`fetcher/fetcher.go`参照）。

## コード設計の注意点

**イベント同一性の比較**
`event.Event.Equals()` はSummary・Start・End の3フィールドで同一性を判断します。Google Calendar のイベントIDでは比較しません（`event/event.go`）。

**差分ロジック**
`event.Events.Diff(another)` は集合演算で実装されています。`another` に存在して `es` にないものが「追加すべきイベント」、`es` に存在して `another` にないものが「削除すべきイベント」になります（`event/events.go`）。

**Google Calendar 認証**
Application Default Credentials（ADC）を使用します。ローカルで動かす場合は `gcloud auth application-default login` が必要です。

## ローカル実行

カレンダーのイベントを全削除する場合：

```bash
CALENDAR_ID=xxx@group.calendar.google.com go run ./cmd/reset/main.go
```

## デプロイ

GCP Cloud Buildを使用します。`cloudbuild.yaml` に定義された手順でCloud Functionsへデプロイされます。

- リージョン: `asia-northeast1`
- ランタイム: `go125`
- トリガー: Pub/Subトピック `fsw-calendar`

## テスト方針

- 全パッケージにユニットテストを記述する
- `go test -race` でデータ競合チェックを行う
- `fetcher` はモックHTTPクライアント（`fetcher.Client`インターフェース）を使ってテストする
- Google Calendar APIは結合テストなし（実APIを直接テストしない）
- カバレッジはGitHub Pagesで公開される

## 注意事項

- このカレンダーはFSW公式情報に基づく非公式な利用であり、可用性は保証されない
- スケジュール変更は最大24時間遅れて反映される
