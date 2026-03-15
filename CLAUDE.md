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
- `logger/` - zapベースのロガー初期化
- `notify/slack/` - Slackへのエラー通知
- `pages/` - テストカバレッジレポート（GitHub Pages）
- `cmd/` - ローカル実行用コマンド

## 開発コマンド

```bash
# 依存関係のインストール
make setup

# テスト実行
make test

# カバレッジレポート生成・表示
make coverage
```

## アーキテクチャ

1. Cloud SchedulerがPub/SubトピックにメッセージをPublish
2. `Register`関数がPub/Subメッセージを受信
3. メッセージからコース・クラス・カレンダーIDを取得
4. FSWサイトから該当月のスケジュールをスクレイピング
5. 既存のGoogleカレンダーイベントと差分比較
6. 追加・削除が必要なイベントのみ更新
7. エラー発生時はSlackに通知

## 対応カレンダー

| カテゴリ | 説明 |
|----------|------|
| SS-4 | スーパースポーツ走行枠 |
| T-4  | ツーリング走行枠 |
| NS-4 | ノーマルスポーツ走行枠 |
| S-4  | スポーツ走行枠 |

## 環境変数

設定は`config/`パッケージで管理されています。デプロイ時は`cloudbuild.yaml`および`config/`内の設定を参照してください。

## デプロイ

GCP Cloud Buildを使用してデプロイします。

```bash
# Cloud Buildで自動デプロイ（mainブランチへのpushで実行）
# cloudbuild.yaml を参照
```

## テスト方針

- 全パッケージにユニットテストを記述する
- `go test -race` でデータ競合チェックを行う
- カバレッジはGitHub Pagesで公開される

## 注意事項

- このカレンダーはFSW公式情報に基づく非公式な利用であり、可用性は保証されない
- スケジュール変更は最大24時間遅れて反映される
