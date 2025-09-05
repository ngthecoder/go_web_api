# Go レシピ Web API + Next.js
Go言語のバックエンドAPIとNext.jsフロントエンドを使用したWebアプリケーション

## プロジェクトの目的
- Go言語を使用したWebAPI開発の学習
- 食材とレシピの相互参照システムの構築
- チーム開発の経験を積む

## 技術スタック
- **バックエンド**： Go 1.22.4 + SQLite3
- **フロントエンド**： Next.js 15.5.0 + Typescript + Tailwind.css
- **データベース**: SQLite3（ローカル開発用）

## 前提条件
開始する前に以下のものがインストールされていることを確認する
1. Go 1.22以上
2. Node.js 18以上
3. Git
4. VSCodeなどのテキストエディタ
5. Docker Desktop（本番用、任意）

## はじめ方
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

## システムアーキテクチャ
### データベース設計
#### テーブル構造
ingredients（食材テーブル）
- id (INTEGER)
- name (TEXT)
- category (TEXT)
- calories_per_100g (INTEGER)
- description (TEXT)
- PRIMARY KEY: id

recipes（レシピテーブル）
- id (INTEGER)
- name (TEXT)
- category (TEXT)
- prep_time_minutes (INTEGER)
- cook_time_minutes (INTEGER)
- servings (INTEGER)
- difficulty (TEXT)
- instructions (TEXT)
- description (TEXT)
- PRIMARY KEY: id

recipe_ingredients（レシピ-食材関連テーブル）
- recipe_id (INTEGER)
- ingredient_id (INTEGER)
- quantity (REAL)
- unit (TEXT)
- notes (TEXT)
- PRIMARY KEY: recipe_id, ingredient_id
- FOREIGN KEY: recipe_id, ingredient_id

#### recipe_ingredientsテーブルの必要性について
##### 悪い例
recipesテープル：
```
| id | name | ingredients_list |
|----|------|------------------|
| 1. | トマトライス | "トマト、玉ねぎ、米" |
| 2. | トマトパスタ | "トマト、玉ねぎ、ニンニク" |
```

問題点：
- 検索が困難（例えばトマトだけを使うレシピの検索ができない）
- 分量情報を保存できない

##### 良い例
ingredientsテーブル：
```
| id | name |
|----|------|
| 1 | トマト |
| 2 | 玉ねぎ |
```

recipesテープル：
```
| id | name |
|----|------|
| 1 | トマトライス |
| 2 | トマトパスタ |
```

recipe_ingredientsテーブル：
```
| recipe_id | ingredients_id | quantity | unit | notes |
|-----------|----------------|----------|------|-------|
| 1 | 1 | 2 | 個 | 角切り |
| 1 | 2 | 1 | 個 | みじん切り |
| 2 | 1 | 3 | 個 | スライス |
| 2 | 2 | 1 | 個 | みじん切り |
```

利点：
- 効率的な検索が可能
- 詳細情報の保存
- 柔軟性

#### データベース関係図
```
ingredients (食材)     recipe_ingredients (中間)     recipes (レシピ)
┌──────────────┐      ┌─────────────────────┐      ┌─────────────────┐
│ id (PK)      │◄─────┤ ingredient_id (FK)  │      │ id (PK)         │
│ name         │      │ recipe_id (FK)      ├─────►│ name            │
│ category     │      │ quantity            │      │ category        │
│ calories     │      │ unit                │      │ prep_time       │
│ description  │      │ notes               │      │ cook_time       │
└──────────────┘      └─────────────────────┘      │ servings        │
                                                   │ difficulty      │
                                                   │ instructions    │
                                                   │ description     │
                                                   └─────────────────┘
```

### APIエンドポイント設計
#### エンドポイント一覧
| HTTP Method | エンドポイント | 説明 | 必須パラメータ | オプションパラメータ |
|-------------|---------------|------|-------------|-------------|
| GET | `/ingredients` | 食材一覧取得・検索・フィルタリング | - | `search`, `category`, `sort`, `order`, `page`, `limit` |
| GET | `/ingredients/{id}` | 食材詳細＋関連レシピ取得 | `id` | - |
| GET | `/recipes` | レシピ一覧取得・検索・フィルタリング | - | `search`, `category`, `max_time`, `difficulty`, `sort`, `order`, `page`, `limit` |
| GET | `/recipes/find-by-ingredients` | 手持ち食材からレシピ検索 | `ingredients` | `match_type`, `page`, `limit` |
| GET | `/recipes/{id}` | レシピ詳細＋使用食材取得 | `id` | - |
| GET | `/recipes/shopping-list/{id}` | レシピの買い物リスト生成 | `id` | `have_ingredients` |
| GET | `/categories` | カテゴリ統計データ取得 | - | - |
| GET | `/stats` | 全体統計情報取得 | - | - |

#### 詳細仕様
**1: GET /api/recipes/find-by-ingredients**
```bash
GET /api/recipes/find-by-ingredients?ingredients=2,26&match_type=partial&limit=1
```
```json
{
  "matched_recipes": [
    {
      "id": 5,
      "name": "オムライス",
      "category": "昼食",
      "prep_time_minutes": 15,
      "cook_time_minutes": 20,
      "servings": 2,
      "difficulty": "medium",
      "instructions": "1. 玉ねぎを炒める\n2. ご飯を加えてチャーハンを作る\n3. 卵を溶いて薄焼き卵を作る\n4. チャーハンを包む",
      "description": "みんな大好きオムライス",
      "matched_ingredients_count": 2,
      "total_ingredients_count": 6,
      "match_score": 0.33333334
    }
  ]
}
```

**2: GET /api/recipes/shopping-list/{id}**
```bash
# 例：夏野菜カレー(id=5)の買い物リスト、トマトとナスは持っている
GET /api/recipes/shopping-list/5?have_ingredients=1,4
```
```json
{
  "recipe": {
    "id": 5,
    "name": "夏野菜カレー",
    "servings": 4
  },
  "need_to_buy": [
    {
      "id": 6,
      "name": "ジャガイモ",
      "quantity": 3,
      "unit": "個",
      "estimated_price": "¥150",
      "category": "野菜"
    },
    {
      "id": 8,
      "name": "カレールー",
      "quantity": 1,
      "unit": "箱",
      "estimated_price": "¥200",
      "category": "調味料"
    }
  ],
  "already_have": [
    {"id": 1, "name": "トマト", "quantity": 2, "unit": "個"},
    {"id": 4, "name": "ナス", "quantity": 1, "unit": "個"}
  ],
  "total_estimated_cost": "¥350"
}
```

**3: GET /api/categories**
```json
{
  "ingredient_categories": {
    "野菜": 15,
    "タンパク質": 8,
    "穀物": 5,
    "乳製品": 4
  },
  "recipe_categories": {
    "朝食": 12,
    "昼食": 18,
    "夕食": 25,
    "おやつ": 8
  }
}
```

**4: GET /api/stats**
```json
{
  "total_ingredients": 32,
  "total_recipes": 63,
  "avg_prep_time": 15.5,
  "avg_cook_time": 22.3,
  "difficulty_distribution": {
    "easy": 45,
    "medium": 15,
    "hard": 3
  }
}
```

**5: GET /api/ingredients**
```bash
GET /api/ingredients?search=トマト&category=野菜&sort=calories&order=desc&page=1&limit=10
```
```json
{
  "has_next": false,
  "ingredients": [
    {
      "id": 1,
      "name": "トマト",
      "category": "野菜",
      "calories_per_100g": 18,
      "description": "新鮮な赤いトマト"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total": 1,
  "total_pages": 1
}
```

**6: GET /api/ingredients/{id}**
```json
{
  "ingredient": {
    "id": 1,
    "name": "トマト",
    "category": "野菜",
    "calories_per_100g": 18,
    "description": "新鮮な赤いトマト"
  },
  "recipes": [
    {
      "id": 1,
      "name": "トマトライス",
      "category": "夕食",
      "prep_time_minutes": 10,
      "cook_time_minutes": 25,
      "servings": 4,
      "difficulty": "easy",
      "instructions": "1. フライパンで油を熱する...",
      "description": "シンプルで美味しいトマトライス"
    }
  ]
}
```

**7: GET /api/recipes**
```bash
GET /api/recipes?search=カレー&category=夕食&difficulty=medium&max_time=60&sort=total_time&order=asc&page=1&limit=5
```
```json
{
  "has_next": false,
  "page": 1,
  "page_size": 5,
  "recipes": [
    {
      "id": 13,
      "name": "夏野菜カレー",
      "category": "夕食",
      "prep_time_minutes": 20,
      "cook_time_minutes": 30,
      "servings": 4,
      "difficulty": "medium",
      "instructions": "1. 野菜をカットする\n2. 鍋で野菜を炒める\n3. 水を加えて煮込む\n4. カレールーを溶かし入れる\n5. さらに煮込んで完成",
      "description": "夏野菜たっぷりのヘルシーカレー"
    },
    {
      "id": 14,
      "name": "チキンカレー",
      "category": "夕食",
      "prep_time_minutes": 15,
      "cook_time_minutes": 45,
      "servings": 4,
      "difficulty": "medium",
      "instructions": "1. 鶏肉をカットする\n2. 玉ねぎを炒める\n3. 鶏肉を加えて炒める\n4. 水を加えて煮込む\n5. カレールーを加える",
      "description": "本格的なチキンカレー"
    }
  ],
  "total": 2,
  "total_pages": 1
}
```

**8: GET /api/recipes/{id}**
```json
{
  "recipe": {
    "id": 1,
    "name": "トマトライス",
    "category": "夕食", 
    "prep_time_minutes": 10,
    "cook_time_minutes": 25,
    "servings": 4,
    "difficulty": "easy",
    "instructions": "1. フライパンで油を熱する...",
    "description": "シンプルで美味しいトマトライス"
  },
  "ingredients": [
    {
      "ingredient_id": 1,
      "name": "トマト",
      "quantity": 2,
      "unit": "個",
      "notes": "角切り"
    },
    {
      "ingredient_id": 2,
      "name": "玉ねぎ",
      "quantity": 1,
      "unit": "個", 
      "notes": "みじん切り"
    }
  ]
}
```

## フロントエンド設計（Next.js）
### ページ詳細
#### メインページ（/）
- 食材とレシピの切り替えタブ
- 検索・フィルタリング機能
- ページネーション
- カード形式でのデータ表示

#### 食材詳細ページ（/ingredients/[id]）
- 食材の詳細情報表示（カロリー、カテゴリ、説明）
- この食材を使用するレシピ一覧
- レシピカードをクリックでレシピ詳細に移動

#### レシピ詳細ページ（/recipes/[id]）
- レシピの詳細情報（調理時間、人数、難易度、手順）
- 必要な食材一覧（分量・単位付き）
- 食材をクリックで食材詳細に移動
- 買い物リスト生成ボタン

### ユーザーフロー
#### 流れ1：手持ち食材からレシピ検索
1. メインページで "手持ち食材で検索" タブを選択
2. 食材リストから持っている食材をチェック　(例: トマト、玉ねぎ、ナス)
3. "レシピを検索" ボタンをクリック
4. マッチしたレシピ一覧表示 (スコア順)
5. 気になるレシピをクリック
6. レシピ詳細ページで手順と全材料を確認
7. "買い物リスト作成" ボタンで足りない材料を確認

#### 流れ2:：レシピから食材詳細への移動
1. レシピ詳細ページで材料一覧を表示
2. 気になる食材（例: "ナス"）をクリック
3. ナスの詳細ページへ移動
4. 他のレシピに興味があればクリックして移動

### プロジェクト構成
```
frontend/
├── app/
│   ├── page.tsx              # メインページ
│   ├── ingredients/
│   │   └── [id]/page.tsx     # 食材詳細ページ
│   ├── recipes/
│   │   └── [id]/page.tsx     # レシピ詳細ページ
│   └── globals.css           # Tailwind CSS
├── components/
│   ├── IngredientCard.tsx    # 食材カードコンポーネント
│   ├── RecipeCard.tsx        # レシピカードコンポーネント
│   ├── SearchBar.tsx         # 検索バーコンポーネント
│   └── Pagination.tsx        # ページネーションコンポーネント
└── lib/
    └── api.ts                # API通信用関数
```

## 機能使用
### ユーザーストーリー
1. 食材一覧閲覧: ユーザーは食材一覧を検索・フィルタリングして閲覧できる
2. 食材詳細表示: ユーザーは食材をクリックして詳細とそれを使うレシピを見ることができる
3. レシピ一覧閲覧: ユーザーはレシピ一覧を検索・フィルタリングして閲覧できる
4. レシピ詳細表示: ユーザーはレシピをクリックして詳細と必要な食材を見ることができる
5. 手持ち食材レシピ検索: ユーザーは冷蔵庫にある食材から作れるレシピを検索できる
6. 買い物リスト生成: ユーザーはレシピを選んで、足りない食材の買い物リストを作成できる
7. 双方向ナビゲーション: 食材→レシピ→食材のサイクル移動ができる
8. 統計情報表示: ユーザーはカテゴリ別統計や全体統計を閲覧できる
9. ページ分け表示: 大量データを効率的にページ分けして閲覧できる

### 使用例
**シナリオ1：冷蔵庫の食材で料理を作りたい**
1. 手持ち食材（トマト、ナス、ジャガイモ）を選択
2. システムが作れるレシピを提案（マッチ度付き）
3. レシピを選択して詳細を確認

**シナリオ2：作りたい料理の買い物リストを作成**
1. 作りたいレシピ（例：夏野菜カレー）を選択
2. 手持ち食材を指定
3. システムが買い物リストを自動生成（推定価格付き）

## チーム開発フロー
### バックエンドタスク
- データベーススキーマ作成
- 基本GETエンドポイント実装
- 手持ち食材レシピ検索機能実装
- 買い物リスト生成機能実装
- 検索・フィルタリング機能実装
- ページネーション実装
- データ集計エンドポイント追加
- マッチングスコア算出ロジック実装
- エラーハンドリング強化
- パフォーマンス最適化（インデックス追加等）
- APIドキュメンテーション専用サーバー実装

### フロントエンドタスク（Next.js + TypeScript）
- Next.js プロジェクトセットアップ
- API通信用関数実装（lib/api.ts）
- 手持ち食材選択UI実装
- レシピ提案表示機能
- 買い物リスト表示機能
- 検索・フィルタUI実装（SearchBar.tsx）
- ページネーション UI実装（Pagination.tsx）
- 食材・レシピカードコンポーネント実装
- レスポンシブデザイン（Tailwind CSS）
- 双方向ナビゲーション機能
- ローディング・エラー状態の処理

#### ドキュメンテーションタスク
- API仕様書作成（静的HTMLサイト）
- エンドポイント使用例・サンプルコード作成
- 開発者向けAPIガイド作成
- レスポンス形式説明書作成


## トラブルシューティング
### よくある問題
**1. データベースファイルが見つからない**
```bash
# 解決策：バックエンドディレクトリでサーバーを起動する
cd backend
go run main.go
```

**2. CORS エラー**
```bash
# 現象：フロントエンドからAPIを呼び出せない  
# 解決策：main.goでCORS設定を確認する
```

**3. SQLite3 ドライバーエラー**
```bash
# 解決策：CGO を有効にしてビルドする
CGO_ENABLED=1 go run main.go
```

# 本番環境（未完了）
```bash
docker compose up --build
```