// app/ingredients/[id]/page.tsx
// このファイルはNext.jsのApp Routerで動作します。

// APIから返される各レシピオブジェクトの型を定義
type Recipe = {
  id: number;
  name: string;
  category: string;
  prep_time_minutes: number;
  cook_time_minutes: number;
  servings: number;
  difficulty: string;
  instructions: string;
  description: string;
};

// APIから返されるトップレベルのデータ構造の型を定義
type IngredientDetails = {
  ingredient: {
    id: number;
    name: string;
    category: string;
    calories_per_100g: number;
    description: string;
  };
  recipes: Recipe[];
};

// APIから特定の材料の詳細と関連レシピを取得する関数
async function getIngredientAndRecipes(id: string): Promise<IngredientDetails> {
  try {
    // キャッシュを無効化するためにno-storeを使用 (開発中は重要)
    const res = await fetch(`http://localhost:8000/api/ingredients/${id}`, { cache: 'no-store' });
    
    if (!res.ok) {
      // 特定のIDの食材が見つからなかった場合 (例: 404)
      if (res.status === 404) {
        throw new Error('指定された食材が見つかりませんでした。');
      }
      throw new Error(`APIからのデータ取得に失敗しました。ステータスコード: ${res.status}`);
    }
    
    return res.json();
  } catch (error) {
    console.error('API呼び出し中にエラーが発生しました:', error);
    throw new Error('APIの呼び出しに失敗しました。ネットワーク接続を確認してください。');
  }
}

// 動的ルーティング用のサーバーコンポーネント
export default async function IngredientDetail({ 
  params 
}: { 
  params: { id: string } 
}) {
  let data: IngredientDetails | null = null;
  let error = null;

  const { id } = params;

  // IDがない場合の早期リターン (このルートにアクセスされることは稀だが念のため)
  if (!id) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-red-50 text-red-700 p-8">
        <div className="text-center p-8 bg-white rounded-xl shadow-lg border border-red-200">
          <h2 className="text-2xl font-bold">エラー</h2>
          <p className="mt-4 text-lg">材料のIDが指定されていません。</p>
        </div>
      </div>
    );
  }

  try {
    data = await getIngredientAndRecipes(id);
  } catch (err) {
    error = (err as Error).message;
  }

  if (error) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-red-50 text-red-700 p-8">
        <div className="text-center p-8 bg-white rounded-xl shadow-lg border border-red-200">
          <h2 className="text-2xl font-bold">エラー</h2>
          <p className="mt-4 text-lg">{error}</p>
        </div>
      </div>
    );
  }
  
  // dataがnullでないことを保証
  if (!data || !data.ingredient) { // ingredientが存在しない場合も考慮
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-50 text-gray-700 p-8">
        <p className="text-xl">食材の詳細情報が見つかりませんでした。</p>
      </div>
    );
  }

  const { ingredient, recipes } = data;

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="container mx-auto max-w-3xl">
        {/* 食材の詳細カード */}
        <div className="bg-white p-8 rounded-xl shadow-2xl border border-gray-200 mb-12">
          <div className="flex items-baseline mb-4">
            <span className="bg-blue-100 text-blue-800 text-sm font-semibold px-3 py-1 rounded-full mr-3">
              {ingredient.category}
            </span>
            <h1 className="text-4xl font-extrabold text-gray-900 leading-tight">
              {ingredient.name}
            </h1>
          </div>
          <p className="mt-4 text-gray-700 text-lg leading-relaxed">
            {ingredient.description}
          </p>
          <div className="mt-6 pt-4 border-t border-gray-100">
            <p className="text-gray-600 text-md">
              <span className="font-semibold text-gray-800">100gあたりのカロリー:</span> {ingredient.calories_per_100g} kcal
            </p>
            <p className="text-gray-400 text-sm mt-2">
              食材ID: {ingredient.id}
            </p>
          </div>
        </div>
        
        {/* 関連レシピのリスト */}
        {recipes && recipes.length > 0 && (
          <div className="mt-8">
            <h2 className="text-3xl font-bold text-gray-900 mb-6 border-b-2 border-blue-500 pb-2">この食材を使ったレシピ</h2>
            <ul className="space-y-6">
              {recipes.map((recipe: Recipe) => (
                <li key={recipe.id} className="bg-white p-6 rounded-xl shadow-lg border border-gray-100 hover:shadow-xl transition-shadow duration-300">
                  <h3 className="text-2xl font-semibold text-gray-800 mb-2">{recipe.name}</h3>
                  <div className="flex items-center text-sm text-gray-500 mb-3 space-x-3">
                    <span className="bg-green-100 text-green-800 px-2 py-0.5 rounded-full font-medium">
                      {recipe.category}
                    </span>
                    <span>🕒 {recipe.prep_time_minutes + recipe.cook_time_minutes}分</span>
                    <span>🍴 {recipe.servings}人前</span>
                    <span className={`
                      ${recipe.difficulty === 'easy' ? 'text-green-600' : ''}
                      ${recipe.difficulty === 'medium' ? 'text-orange-600' : ''}
                      ${recipe.difficulty === 'hard' ? 'text-red-600' : ''}
                      font-semibold
                    `}>
                      難易度: {recipe.difficulty === 'easy' ? 'かんたん' : recipe.difficulty === 'medium' ? 'ふつう' : 'むずかしい'}
                    </span>
                  </div>
                  <p className="mt-2 text-gray-600 leading-relaxed">{recipe.description}</p>
                  <details className="mt-4 text-gray-700 cursor-pointer">
                    <summary className="font-semibold text-blue-700 hover:text-blue-800">作り方を見る</summary>
                    <p className="mt-3 whitespace-pre-line text-sm bg-gray-50 p-4 rounded-lg border border-gray-100">
                      {recipe.instructions}
                    </p>
                  </details>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}