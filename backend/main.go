package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Ingredient struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Calories    int    `json:"calories_per_100g"`
	Description string `json:"description"`
}
type Recipe struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Category        string `json:"category"`
	PrepTimeMinutes int    `json:"prep_time_minutes"`
	CookTimeMinutes int    `json:"cook_time_minutes"`
	Servings        int    `json:"servings"`
	Difficulty      string `json:"difficulty"`
	Instructions    string `json:"instructions"`
	Description     string `json:"description"`
}

type RecipeIngredient struct {
	RecipeID     int     `json:"recipe_id"`
	IngredientID int     `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	Notes        string  `json:"notes"`
}

type RecipeWithIngredients struct {
	Recipe      Recipe                   `json:"recipe"`
	Ingredients []IngredientWithQuantity `json:"ingredients"`
}

type IngredientWithQuantity struct {
	IngredientID int     `json:"ingredient_id"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	Notes        string  `json:"notes"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./foods.db")
	if err != nil {
		log.Fatal(err)
	}

	createIngredientsTable := `
		CREATE TABLE IF NOT EXISTS ingredients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			calories_per_100g INTEGER NOT NULL,
			description TEXT
		);
	`

	createRecipesTable := `
		CREATE TABLE IF NOT EXISTS recipes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			prep_time_minutes INTEGER NOT NULL,
			cook_time_minutes INTEGER NOT NULL,
			servings INTEGER NOT NULL,
			difficulty TEXT NOT NULL,
			instructions TEXT NOT NULL,
			description TEXT
		);
	`

	createRecipesIngredientsTable := `
		CREATE TABLE IF NOT EXISTS recipe_ingredients (
			recipe_id INTEGER NOT NULL,
			ingredient_id INTEGER NOT NULL,
			quantity REAL NOT NULL,
			unit TEXT NOT NULL,
			notes TEXT,
			PRIMARY KEY (recipe_id, ingredient_id),
			FOREIGN KEY (recipe_id) REFERENCES recipes (id),
			FOREIGN KEY (ingredient_id) REFERENCES ingredients (id)
		)
	`

	_, err = db.Exec(createIngredientsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createRecipesTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createRecipesIngredientsTable)
	if err != nil {
		log.Fatal(err)
	}
}

func populateTestData() {
	db.Exec("DELETE FROM recipe_ingredients")
	db.Exec("DELETE FROM recipes")
	db.Exec("DELETE FROM ingredients")

	ingredientsData := []struct {
		name        string
		category    string
		calories    int
		description string
	}{
		{"トマト", "野菜", 18, "新鮮な赤いトマト"},
		{"玉ねぎ", "野菜", 37, "甘味のある玉ねぎ"},
		{"ニンニク", "野菜", 134, "香り豊かなニンニク"},
		{"ナス", "野菜", 22, "紫色の美味しいナス"},
		{"ジャガイモ", "野菜", 77, "ホクホクのジャガイモ"},
		{"人参", "野菜", 39, "オレンジ色の甘い人参"},
		{"レタス", "野菜", 12, "新鮮なレタス"},
		{"キュウリ", "野菜", 14, "みずみずしいキュウリ"},
		{"ピーマン", "野菜", 22, "緑のピーマン"},
		{"もやし", "野菜", 14, "シャキシャキもやし"},
		{"キャベツ", "野菜", 23, "甘いキャベツ"},
		{"ブロッコリー", "野菜", 33, "栄養豊富なブロッコリー"},
		{"ほうれん草", "野菜", 20, "鉄分豊富なほうれん草"},
		{"白菜", "野菜", 14, "みずみずしい白菜"},
		{"大根", "野菜", 18, "辛味のある大根"},
		{"かぼちゃ", "野菜", 49, "甘いかぼちゃ"},
		{"きのこ", "野菜", 22, "うまみたっぷりしいたけ"},
		{"えのき", "野菜", 22, "食感の良いえのき"},
		{"しめじ", "野菜", 18, "香り豊かなしめじ"},
		{"アスパラガス", "野菜", 22, "春の味覚アスパラガス"},
		{"セロリ", "野菜", 15, "シャキシャキセロリ"},
		{"生姜", "野菜", 30, "辛みのある生姜"},
		{"長ねぎ", "野菜", 28, "薬味に最適な長ねぎ"},
		{"小松菜", "野菜", 14, "カルシウム豊富な小松菜"},
		{"水菜", "野菜", 23, "シャキシャキ水菜"},

		{"米", "穀物", 358, "日本のお米"},
		{"パスタ", "穀物", 371, "イタリアンパスタ"},
		{"パン", "穀物", 264, "食パン"},
		{"うどん", "穀物", 270, "讃岐うどん"},
		{"そば", "穀物", 274, "日本そば"},
		{"ラーメン", "穀物", 281, "中華麺"},
		{"オートミール", "穀物", 380, "健康的なオートミール"},
		{"小麦粉", "穀物", 368, "薄力粉"},
		{"パン粉", "穀物", 373, "サクサクパン粉"},

		{"鶏肉", "タンパク質", 200, "新鮮な鶏胸肉"},
		{"豚肉", "タンパク質", 263, "柔らかい豚ロース"},
		{"牛肉", "タンパク質", 250, "上質な牛肉"},
		{"卵", "タンパク質", 151, "新鮮な鶏卵"},
		{"鮭", "タンパク質", 139, "脂の乗った鮭"},
		{"マグロ", "タンパク質", 125, "赤身のマグロ"},
		{"エビ", "タンパク質", 83, "プリプリエビ"},
		{"イカ", "タンパク質", 83, "新鮮なイカ"},
		{"豆腐", "タンパク質", 56, "絹ごし豆腐"},
		{"納豆", "タンパク質", 200, "栄養豊富な納豆"},
		{"ひき肉", "タンパク質", 224, "合いびき肉"},
		{"ソーセージ", "タンパク質", 321, "ジューシーソーセージ"},
		{"ハム", "タンパク質", 196, "スモークハム"},
		{"ベーコン", "タンパク質", 405, "カリカリベーコン"},
		{"鶏もも肉", "タンパク質", 253, "ジューシーな鶏もも肉"},

		{"牛乳", "乳製品", 67, "新鮮な牛乳"},
		{"チーズ", "乳製品", 113, "濃厚なチーズ"},
		{"バター", "乳製品", 745, "無塩バター"},
		{"ヨーグルト", "乳製品", 62, "プレーンヨーグルト"},
		{"生クリーム", "乳製品", 433, "濃厚な生クリーム"},
		{"クリームチーズ", "乳製品", 346, "なめらかクリームチーズ"},
		{"モッツァレラ", "乳製品", 280, "イタリアンモッツァレラ"},

		{"オリーブオイル", "調味料", 884, "エクストラバージンオリーブオイル"},
		{"醤油", "調味料", 71, "日本の醤油"},
		{"塩", "調味料", 0, "海塩"},
		{"砂糖", "調味料", 387, "白砂糖"},
		{"カレールー", "調味料", 512, "カレーのルー"},
		{"みそ", "調味料", 217, "赤味噌"},
		{"みりん", "調味料", 241, "本みりん"},
		{"酢", "調味料", 25, "米酢"},
		{"料理酒", "調味料", 103, "日本酒"},
		{"ごま油", "調味料", 921, "香ばしいごま油"},
		{"マヨネーズ", "調味料", 703, "まろやかマヨネーズ"},
		{"ケチャップ", "調味料", 119, "トマトケチャップ"},
		{"ウスターソース", "調味料", 117, "濃厚ウスターソース"},
		{"コショウ", "調味料", 251, "黒胡椒"},
		{"唐辛子", "調味料", 419, "一味唐辛子"},
		{"わさび", "調味料", 265, "本わさび"},
		{"からし", "調味料", 336, "和からし"},
		{"バジル", "調味料", 22, "フレッシュバジル"},
		{"パセリ", "調味料", 43, "イタリアンパセリ"},
		{"ローリエ", "調味料", 313, "月桂樹の葉"},

		{"りんご", "果物", 54, "甘酸っぱいりんご"},
		{"バナナ", "果物", 86, "栄養豊富なバナナ"},
		{"レモン", "果物", 54, "酸っぱいレモン"},
		{"オレンジ", "果物", 39, "ビタミンC豊富なオレンジ"},
		{"いちご", "果物", 34, "甘いいちご"},
		{"ブルーベリー", "果物", 49, "抗酸化作用のあるブルーベリー"},
		{"アボカド", "果物", 187, "濃厚なアボカド"},

		{"はちみつ", "その他", 294, "天然はちみつ"},
		{"のり", "その他", 188, "海苔"},
		{"ごま", "その他", 578, "白ごま"},
		{"片栗粉", "その他", 330, "とろみ付け用"},
		{"コンソメ", "その他", 235, "洋風だしの素"},
		{"だしの素", "その他", 276, "和風だしの素"},
	}

	for _, ing := range ingredientsData {
		_, err := db.Exec("INSERT INTO ingredients (name, category, calories_per_100g, description) VALUES (?, ?, ?, ?)",
			ing.name, ing.category, ing.calories, ing.description)
		if err != nil {
			log.Printf("Error inserting ingredient %s: %v", ing.name, err)
		}
	}

	recipesData := []struct {
		name         string
		category     string
		prepTime     int
		cookTime     int
		servings     int
		difficulty   string
		instructions string
		description  string
	}{
		{
			"フレンチトースト", "朝食", 5, 10, 2, "easy",
			"1. 卵と牛乳を混ぜる\n2. パンを浸す\n3. フライパンで焼く\n4. はちみつをかけて完成",
			"甘くて美味しい朝食",
		},
		{
			"スクランブルエッグ", "朝食", 3, 5, 1, "easy",
			"1. 卵を溶く\n2. バターでゆっくり炒める\n3. 塩コショウで味付け",
			"ふわふわのスクランブルエッグ",
		},
		{
			"オートミール", "朝食", 2, 5, 1, "easy",
			"1. オートミールに牛乳を加える\n2. 電子レンジで加熱\n3. フルーツを添える",
			"健康的な朝食オートミール",
		},
		{
			"パンケーキ", "朝食", 10, 15, 3, "medium",
			"1. 小麦粉、卵、牛乳を混ぜる\n2. フライパンで焼く\n3. はちみつをかける",
			"ふわふわパンケーキ",
		},

		{
			"オムライス", "昼食", 15, 20, 2, "medium",
			"1. 玉ねぎを炒める\n2. ご飯を加えてチャーハンを作る\n3. 卵を溶いて薄焼き卵を作る\n4. チャーハンを包む",
			"みんな大好きオムライス",
		},
		{
			"カルボナーラ", "昼食", 10, 15, 2, "medium",
			"1. パスタを茹でる\n2. ベーコンを炒める\n3. 卵とチーズを混ぜる\n4. 全て和える",
			"濃厚なカルボナーラ",
		},
		{
			"チャーハン", "昼食", 10, 15, 2, "easy",
			"1. 卵でご飯を炒める\n2. 具材を加える\n3. 醤油で味付け",
			"パラパラチャーハン",
		},
		{
			"焼きそば", "昼食", 10, 10, 2, "easy",
			"1. 麺を茹でる\n2. 野菜と炒める\n3. ソースで味付け",
			"定番の焼きそば",
		},
		{
			"サンドイッチ", "昼食", 10, 0, 1, "easy",
			"1. パンにバターを塗る\n2. 具材を挟む\n3. 食べやすく切る",
			"簡単サンドイッチ",
		},
		{
			"親子丼", "昼食", 10, 15, 2, "medium",
			"1. 鶏肉を煮る\n2. 卵を溶き入れる\n3. ご飯に乗せる",
			"やさしい味の親子丼",
		},

		{
			"トマトライス", "夕食", 10, 25, 4, "easy",
			"1. 玉ねぎをみじん切りにする\n2. フライパンで油を熱し、玉ねぎを炒める\n3. トマトを角切りにして加える\n4. 米を加えて炒める\n5. 水を加えて煮る",
			"シンプルで美味しいトマトライス",
		},
		{
			"トマトパスタ", "夕食", 15, 20, 2, "easy",
			"1. パスタを茹でる\n2. ニンニクをみじん切りにして炒める\n3. トマトを加えて煮込む\n4. 茹でたパスタと和える",
			"基本のトマトパスタ",
		},
		{
			"夏野菜カレー", "夕食", 20, 30, 4, "medium",
			"1. 野菜をカットする\n2. 鍋で野菜を炒める\n3. 水を加えて煮込む\n4. カレールーを溶かし入れる\n5. さらに煮込んで完成",
			"夏野菜たっぷりのヘルシーカレー",
		},
		{
			"チキンカレー", "夕食", 15, 45, 4, "medium",
			"1. 鶏肉をカットする\n2. 玉ねぎを炒める\n3. 鶏肉を加えて炒める\n4. 水を加えて煮込む\n5. カレールーを加える",
			"本格的なチキンカレー",
		},
		{
			"ハンバーグ", "夕食", 20, 25, 4, "medium",
			"1. ひき肉と玉ねぎを混ぜる\n2. 成形する\n3. フライパンで焼く\n4. ソースを作る",
			"ジューシーハンバーグ",
		},
		{
			"生姜焼き", "夕食", 10, 15, 2, "easy",
			"1. 豚肉を切る\n2. 生姜だれを作る\n3. 肉を焼いてたれを絡める",
			"ご飯が進む生姜焼き",
		},
		{
			"鮭の塩焼き", "夕食", 5, 15, 2, "easy",
			"1. 鮭に塩を振る\n2. グリルで焼く\n3. レモンを添える",
			"シンプルな鮭の塩焼き",
		},
		{
			"麻婆豆腐", "夕食", 15, 20, 3, "medium",
			"1. ひき肉を炒める\n2. 豆腐を加える\n3. 調味料で味付け\n4. 片栗粉でとろみをつける",
			"辛くて美味しい麻婆豆腐",
		},
		{
			"エビチリ", "夕食", 20, 15, 3, "hard",
			"1. エビの下処理\n2. 衣をつけて揚げる\n3. チリソースを作る\n4. エビと和える",
			"プリプリエビチリ",
		},
		{
			"肉じゃが", "夕食", 15, 30, 4, "medium",
			"1. 具材を切る\n2. 油で炒める\n3. だしと調味料で煮込む",
			"家庭的な肉じゃが",
		},
		{
			"すき焼き", "夕食", 20, 25, 4, "medium",
			"1. 牛肉と野菜を準備\n2. 鍋で煮る\n3. 生卵につけて食べる",
			"贅沢なすき焼き",
		},
		{
			"天ぷら", "夕食", 25, 20, 4, "hard",
			"1. 衣を作る\n2. 具材に衣をつける\n3. 油で揚げる",
			"サクサク天ぷら",
		},

		{
			"サラダ", "副菜", 10, 0, 2, "easy",
			"1. 野菜を洗って切る\n2. ボウルに盛り付ける\n3. ドレッシングをかける",
			"新鮮野菜のサラダ",
		},
		{
			"味噌汁", "副菜", 5, 10, 4, "easy",
			"1. だしを取る\n2. 具材を煮る\n3. 味噌を溶く",
			"定番の味噌汁",
		},
		{
			"きんぴらごぼう", "副菜", 10, 15, 3, "medium",
			"1. ごぼうと人参を切る\n2. 炒める\n3. 調味料で味付け",
			"シャキシャキきんぴら",
		},
		{
			"クッキー", "おやつ", 30, 20, 12, "medium",
			"1. バターと砂糖を混ぜる\n2. 小麦粉を加える\n3. 成形して焼く",
			"手作りクッキー",
		},
		{
			"プリン", "おやつ", 15, 30, 6, "hard",
			"1. カラメルを作る\n2. プリン液を作る\n3. 蒸し器で蒸す",
			"なめらかプリン",
		},
	}

	for _, rec := range recipesData {
		_, err := db.Exec(`INSERT INTO recipes (name, category, prep_time_minutes, cook_time_minutes, servings, difficulty, instructions, description) 
						  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			rec.name, rec.category, rec.prepTime, rec.cookTime, rec.servings, rec.difficulty, rec.instructions, rec.description)
		if err != nil {
			log.Printf("Error inserting recipe %s: %v", rec.name, err)
		}
	}

	recipeIngredientsData := []struct {
		recipeID     int
		ingredientID int
		quantity     float64
		unit         string
		notes        string
	}{
		{1, 28, 4, "枚", "厚切り"},
		{1, 38, 2, "個", ""},
		{1, 50, 100, "ml", ""},
		{1, 52, 20, "g", ""},
		{1, 84, 2, "大さじ", "仕上げ用"},

		{2, 38, 3, "個", ""},
		{2, 52, 10, "g", ""},
		{2, 59, 1, "小さじ", ""},
		{2, 70, 1, "小さじ", ""},

		{3, 32, 50, "g", ""},
		{3, 50, 150, "ml", ""},
		{3, 77, 1, "個", "スライス"},
		{3, 78, 1, "本", ""},

		{4, 33, 100, "g", ""},
		{4, 38, 1, "個", ""},
		{4, 50, 150, "ml", ""},
		{4, 60, 2, "大さじ", ""},
		{4, 52, 20, "g", ""},

		{5, 2, 1, "個", "みじん切り"},
		{5, 26, 2, "カップ", "冷ご飯"},
		{5, 38, 4, "個", ""},
		{5, 59, 1, "小さじ", ""},
		{5, 52, 10, "g", ""},
		{5, 68, 2, "大さじ", ""},

		{6, 27, 200, "g", ""},
		{6, 48, 80, "g", "スライス"},
		{6, 38, 2, "個", ""},
		{6, 51, 50, "g", "すりおろし"},
		{6, 70, 1, "小さじ", ""},

		{7, 26, 2, "カップ", "冷ご飯"},
		{7, 38, 2, "個", ""},
		{7, 2, 1, "個", "みじん切り"},
		{7, 23, 1, "本", "小口切り"},
		{7, 58, 2, "大さじ", ""},
		{7, 66, 1, "大さじ", ""},

		{8, 31, 2, "玉", ""},
		{8, 11, 100, "g", "細切り"},
		{8, 6, 1, "本", "細切り"},
		{8, 10, 1, "袋", ""},
		{8, 69, 3, "大さじ", ""},

		{9, 28, 4, "枚", ""},
		{9, 52, 20, "g", ""},
		{9, 47, 2, "枚", ""},
		{9, 7, 2, "枚", ""},
		{9, 1, 1, "個", "スライス"},

		{10, 35, 200, "g", "一口大"},
		{10, 2, 1, "個", "スライス"},
		{10, 38, 3, "個", ""},
		{10, 26, 2, "膳", ""},
		{10, 58, 2, "大さじ", ""},
		{10, 63, 2, "大さじ", ""},

		{11, 1, 2, "個", "角切り"},
		{11, 2, 1, "個", "みじん切り"},
		{11, 26, 2, "カップ", ""},
		{11, 57, 2, "大さじ", ""},
		{11, 59, 1, "小さじ", ""},

		{12, 1, 3, "個", "ざく切り"},
		{12, 3, 2, "片", "みじん切り"},
		{12, 27, 200, "g", ""},
		{12, 57, 2, "大さじ", ""},
		{12, 59, 1, "小さじ", ""},

		{13, 1, 2, "個", "乱切り"},
		{13, 4, 1, "本", "乱切り"},
		{13, 5, 2, "個", "乱切り"},
		{13, 6, 1, "本", "乱切り"},
		{13, 61, 1, "箱", ""},
		{13, 35, 200, "g", "一口大"},

		{14, 35, 400, "g", "一口大"},
		{14, 2, 2, "個", "スライス"},
		{14, 61, 1, "箱", ""},
		{14, 5, 1, "個", "乱切り"},
		{14, 6, 1, "本", "乱切り"},

		{15, 45, 300, "g", ""},
		{15, 2, 1, "個", "みじん切り"},
		{15, 38, 1, "個", ""},
		{15, 34, 50, "g", ""},
		{15, 50, 2, "大さじ", ""},

		{16, 36, 300, "g", ""},
		{16, 22, 1, "片", "すりおろし"},
		{16, 58, 2, "大さじ", ""},
		{16, 63, 2, "大さじ", ""},

		{17, 39, 2, "切れ", ""},
		{17, 59, 1, "小さじ", ""},
		{17, 79, 1, "個", "くし切り"},

		{18, 45, 150, "g", ""},
		{18, 43, 1, "丁", "角切り"},
		{18, 62, 2, "大さじ", ""},
		{18, 87, 1, "大さじ", ""},

		{19, 41, 200, "g", ""},
		{19, 1, 2, "個", "角切り"},
		{19, 68, 3, "大さじ", ""},
		{19, 71, 1, "小さじ", ""},

		{20, 36, 200, "g", "薄切り"},
		{20, 5, 3, "個", "乱切り"},
		{20, 6, 1, "本", "乱切り"},
		{20, 2, 1, "個", "くし切り"},
		{20, 58, 3, "大さじ", ""},
		{20, 60, 2, "大さじ", ""},
		{20, 63, 2, "大さじ", ""},
		{20, 89, 1, "個", ""},

		{21, 37, 300, "g", "薄切り"},
		{21, 2, 1, "個", "くし切り"},
		{21, 14, 1, "袋", ""},
		{21, 17, 2, "本", ""},
		{21, 43, 1, "丁", ""},
		{21, 58, 4, "大さじ", ""},
		{21, 60, 3, "大さじ", ""},
		{21, 65, 2, "大さじ", ""},

		{22, 41, 200, "g", ""},
		{22, 4, 1, "本", "輪切り"},
		{22, 20, 1, "本", ""},
		{22, 33, 100, "g", ""},
		{22, 38, 1, "個", "冷水と混ぜる"},

		{23, 7, 4, "枚", "手でちぎる"},
		{23, 8, 1, "本", "スライス"},
		{23, 1, 1, "個", "くし切り"},
		{23, 57, 2, "大さじ", "ドレッシング"},
		{23, 64, 1, "大さじ", ""},

		{24, 43, 1, "丁", "角切り"},
		{24, 24, 2, "本", "小口切り"},
		{24, 62, 2, "大さじ", ""},
		{24, 89, 1, "個", ""},

		{25, 15, 1, "本", "千切り"},
		{25, 6, 1, "本", "千切り"},
		{25, 66, 1, "大さじ", ""},
		{25, 58, 2, "大さじ", ""},
		{25, 63, 1, "大さじ", ""},
		{25, 71, 1, "小さじ", ""},

		{26, 33, 200, "g", ""},
		{26, 52, 100, "g", ""},
		{26, 60, 80, "g", ""},
		{26, 38, 1, "個", ""},

		{27, 38, 3, "個", ""},
		{27, 50, 300, "ml", ""},
		{27, 60, 80, "g", ""},
		{27, 84, 2, "大さじ", ""},
	}

	for _, ri := range recipeIngredientsData {
		_, err := db.Exec(`INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity, unit, notes) 
						  VALUES (?, ?, ?, ?, ?)`,
			ri.recipeID, ri.ingredientID, ri.quantity, ri.unit, ri.notes)
		if err != nil {
			log.Printf("Error inserting recipe_ingredient: %v", err)
		}
	}

	fmt.Println("Test data populated successfully!")
	fmt.Printf("Inserted %d ingredients, %d recipes, and %d recipe-ingredient relationships\n",
		len(ingredientsData), len(recipesData), len(recipeIngredientsData))
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")

	if name == "" {
		name = "匿名"
	}

	response := map[string]string{
		"message": fmt.Sprintf("ようこそ、%sさん！", name),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ingredientsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM ingredients")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ingredients []Ingredient
	for rows.Next() {
		var ingredient Ingredient
		err = rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Category, &ingredient.Calories, &ingredient.Description)
		if err != nil {
			http.Error(w, "Data scanning error", http.StatusInternalServerError)
		}
		ingredients = append(ingredients, ingredient)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ingredients": ingredients,
	})
}

func recipesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM recipes")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.Category, &recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.Servings, &recipe.Difficulty, &recipe.Instructions, &recipe.Description)
		if err != nil {
			http.Error(w, "Data scanning error", http.StatusInternalServerError)
		}
		recipes = append(recipes, recipe)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recipes": recipes,
	})
}

func findRecipesByIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ingredientsParams := query.Get("ingredients")
	matchType := query.Get("match_type")
	limitParams := query.Get("limit")

	if ingredientsParams == "" {
		http.Error(w, "Missing required parameters: ingredients", http.StatusBadRequest)
		return
	}

	if matchType == "" || matchType == "exact" {
		matchType = "partial"
	}

	limit := 10
	if limitParams != "" {
		if l, err := strconv.Atoi(limitParams); err == nil && l > 0 {
			limit = l
		}
	}

	ingredientIDStrings := strings.Split(ingredientsParams, ",")
	ingredientIDs := make([]int, 0, len(ingredientIDStrings))

	for _, idStr := range ingredientIDStrings {
		if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
			ingredientIDs = append(ingredientIDs, id)
		}
	}

	if len(ingredientIDs) == 0 {
		http.Error(w, "Invalid ingredient IDs", http.StatusBadRequest)
		return
	}

	// SQL Query here
	sqlQuery := ""
	args := []interface{}{}

	placeholders := make([]string, 0, len(ingredientIDs))
	for _, ingredientID := range ingredientIDs {
		placeholders = append(placeholders, "?")
		args = append(args, ingredientID)
	}
	args = append(args, limit)

	if matchType == "partial" {
		sqlQuery = fmt.Sprintf(
			`
			SELECT r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description,
			COUNT(ri.ingredient_id) as match_ingredients_count,
			(SELECT COUNT(*) FROM recipe_ingredients WHERE recipe_id = r.id) as total_ingredients_count
			FROM recipes r
			JOIN recipe_ingredients ri on r.id = ri.recipe_id
			WHERE ri.ingredient_id in (%s)
			GROUP BY r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description
			ORDER BY match_ingredients_count DESC, total_ingredients_count ASC
			LIMIT ?
		`, strings.Join(placeholders, ","))
	}

	if matchType == "exact" {
		sqlQuery = fmt.Sprintf(
			`
			SELECT r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description,
			COUNT(ri.ingredient_id) as match_ingredients_count,
			(SELECT COUNT(*) FROM recipe_ingredients WHERE recipe_id = r.id) as total_ingredients_count
			FROM recipes r
			JOIN recipe_ingredients ri on r.id = ri.recipe_id
			WHERE ri.ingredient_id in (%s)
			GROUP BY r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description
			HAVING COUNT(ri.ingredient_id) = (SELECT COUNT(*) FROM recipe_ingredients WHERE recipe_id = r.id)
			ORDER BY match_ingredients_count DESC, total_ingredients_count ASC
			LIMIT ?
		`, strings.Join(placeholders, ","))
	}

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("SQL Error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type MatchedRecipe struct {
		ID                      int     `json:"id"`
		Name                    string  `json:"name"`
		Category                string  `json:"category"`
		PrepTimeMinutes         int     `json:"prep_time_minutes"`
		CookTimeMinutes         int     `json:"cook_time_minutes"`
		Servings                int     `json:"servings"`
		Difficulty              string  `json:"difficulty"`
		Instructions            string  `json:"instructions"`
		Description             string  `json:"description"`
		MatchedIngredientsCount int     `json:"matched_ingredients_count"`
		TotalIngredientsCount   int     `json:"total_ingredients_count"`
		MatchScore              float32 `json:"match_score"`
	}

	matchedRecipes := []MatchedRecipe{}
	for rows.Next() {
		var matchedRecipe MatchedRecipe
		err = rows.Scan(
			&matchedRecipe.ID,
			&matchedRecipe.Name,
			&matchedRecipe.Category,
			&matchedRecipe.PrepTimeMinutes,
			&matchedRecipe.CookTimeMinutes,
			&matchedRecipe.Servings,
			&matchedRecipe.Difficulty,
			&matchedRecipe.Instructions,
			&matchedRecipe.Description,
			&matchedRecipe.MatchedIngredientsCount,
			&matchedRecipe.TotalIngredientsCount,
		)
		matchedRecipe.MatchScore = float32(matchedRecipe.MatchedIngredientsCount) / float32(matchedRecipe.TotalIngredientsCount)
		if err != nil {
			http.Error(w, "Database scanning error", http.StatusInternalServerError)
			return
		}
		matchedRecipes = append(matchedRecipes, matchedRecipe)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"matched_recipes": matchedRecipes,
	})
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	type Stats struct {
		TotalIngredients       int            `json:"total_ingredients"`
		TotalRecipes           int            `json:"total_recipes"`
		AvgPrepTime            float32        `json:"avg_prep_time"`
		AvgCookTime            float32        `json:"avg_cook_time"`
		DifficultyDistribution map[string]int `json:"difficulty_distribution"`
	}

	var stats Stats
	stats.DifficultyDistribution = make(map[string]int)

	err := db.QueryRow(`SELECT COUNT(*) FROM ingredients`).Scan(&stats.TotalIngredients)
	if err != nil {
		http.Error(w, "Database scanning error", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM recipes`).Scan(&stats.TotalRecipes)
	if err != nil {
		http.Error(w, "Database scanning error", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(`SELECT AVG(prep_time_minutes) FROM recipes`).Scan(&stats.AvgPrepTime)
	if err != nil {
		http.Error(w, "Database scanning error", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(`SELECT AVG(cook_time_minutes) FROM recipes`).Scan(&stats.AvgCookTime)
	if err != nil {
		http.Error(w, "Database scanning error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT difficulty, COUNT(*)
		FROM recipes
		GROUP BY difficulty
	`)
	if err != nil {
		log.Print(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var difficulty string
		var count int
		err = rows.Scan(&difficulty, &count)
		if err != nil {
			log.Print(err)
			http.Error(w, "Database scanning error", http.StatusInternalServerError)
			return
		}
		stats.DifficultyDistribution[difficulty] = count
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	initDB()
	// populateTestData()
	defer db.Close()

	fmt.Printf("ポート8000でAPIサーバーを起動\n")

	http.HandleFunc("/api/hello", enableCORS(helloHandler))
	http.HandleFunc("/api/ingredients", enableCORS(ingredientsHandler))
	http.HandleFunc("/api/recipes", enableCORS(recipesHandler))
	http.HandleFunc("/api/stats", enableCORS(statsHandler))
	http.HandleFunc("/api/recipes/find-by-ingredients", enableCORS(findRecipesByIngredientsHandler))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
