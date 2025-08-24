# Go Web API + Next.js
Go言語のバックエンドAPIとNext.jsフロントエンドを使用したWebアプリケーション

# 技術スタック
- バックエンド：Go 1.22.4
- フロントエンド：Next.js 15.5.0 + Typescript + Tailwind.css

# 前提条件
開始する前に以下のものがインストールされていることを確認する
1. Go 1.22以上
2. Node.js 18以上
3. Git
4. VSCodeなどのテキストエディタ
5. Docker Desktop（本番用、任意）

# はじめ方
1. レポジトリをクローン
```bash
git clone git@github.com:ngthecoder/go_web_api.git
cd go_web_api
```
2. バックエンドの起動
```bash
cd backend
go run main.go
```
- http://localhost:8000で起動
3. フロントエンドの起動
```bash
cd frontend
npm install
npm run dev
```
- http://localhost:3000で起動

# APIテスト
ブラウザかcurlなどを使用してテストAPIの動作確認をしてください
- ブラウザ：http://localhost:8000/api/hello?name={任意の名前}
- コマンドライン
```bash
curl "http://localhost:8000/api/hello?name={任意の名前}"
```

# 本番環境（未完了）
```bash
docker compose up --build
```