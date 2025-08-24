# Go Web API
Go言語とDockerを使用したWeb APIプロジェクト

# 前提条件
開始する前に以下のものがインストールされていることを確認する
1. Docker Desktop
2. Git
3. VSCodeなどのテキストエディタ

# はじめ方
1. レポジトリをクローン
```bash
git clone git@github.com:ngthecoder/go_web_api.git
cd go_web_api
```
2. Dockerの動作確認
```bash
docker --version
docker compose -version
```
3. Docker Containerを起動
```bash
docker compose up --build api
```

# テストAPI
ブラウザかcurlなどを使用してテストAPIの動作確認をしてください
- ブラウザ：http://localhost:8000/hello?name={任意の名前}
- コマンドライン
```bash
curl "http://localhost:8000/hello?name={任意の名前}"
```

# 開発ワークフロー
## 開発モード
```bash
docker compose up --build api
```
- ホスト上のポード8000で実行
- コンパイル済みバイナリを使用
- 変更の反映のため`--build`フラッグ

## サービスの停止
```bash
docker compose down
```

## ログの表示
```bash
docker compose logs -f api
```