'use client';
import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';

interface Ingredient {
  id: number;
  name: string;
  category: string;
  calories_per_100g: number;
  description: string;
}

interface Recipe {
  id: number;
  name: string;
  category: string;
  prep_time_minutes: number;
  cook_time_minutes: number;
  servings: number;
  difficulty: string;
  instructions: string;
  description: string;
}

interface IngredientDetail {
  ingredient: Ingredient;
  recipes: Recipe[];
}

export default function IngredientDetailPage() {
  const params = useParams();
  const [ingredientDetail, setIngredientDetail] = useState<IngredientDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (params.id) {
      fetch(`http://localhost:8000/api/ingredients/${params.id}`)
        .then(res => {
          if (!res.ok) {
            throw new Error(`食材が見つかりません (${res.status})`);
          }
          return res.json();
        })
        .then(data => {
            console.log('API Response:', data);
            setIngredientDetail(data);
            setLoading(false);
        })
        .catch(err => {
            console.error('API Error:', err);
            setError(err.message);
            setLoading(false);
        });
    }
  }, [params.id]);

  const getCategoryColor = (category: string) => {
    const colors: Record<string, string> = {
      '野菜': 'bg-green-100 text-green-800',
      '穀物': 'bg-yellow-100 text-yellow-800',
      'タンパク質': 'bg-red-100 text-red-800',
      '乳製品': 'bg-blue-100 text-blue-800',
      '調味料': 'bg-purple-100 text-purple-800',
      '果物': 'bg-pink-100 text-pink-800',
      'その他': 'bg-gray-100 text-gray-800',
    };
    return colors[category] || 'bg-gray-100 text-gray-800';
  };

  if (loading) return <div className="container mx-auto p-4">Loading...</div>;
  if (error) return <div className="container mx-auto p-4 text-red-600">エラー: {error}</div>;
  if (!ingredientDetail) return <div className="container mx-auto p-4">食材が見つかりません</div>;

  const ingredient = ingredientDetail['ingredient'];
  const recipes = ingredientDetail['recipes']

  return (
    <div className="container mx-auto p-4">
      <nav className="mb-6 text-sm text-gray-600">
        <Link href="/" className="hover:text-blue-600">ホーム</Link>
        <span className="mx-2">→</span>
        <Link href="/ingredients" className="hover:text-blue-600">食材一覧</Link>
        <span className="mx-2">→</span>
        <span className="text-gray-900">{ingredient.name}</span>
      </nav>

      <div className="bg-white rounded-lg shadow-sm border p-6 mb-8">
        <div className="flex justify-between items-start mb-4">
          <h1 className="text-4xl font-bold">{ingredient.name}</h1>
          <span className={`px-3 py-1 rounded-full text-sm font-medium ${getCategoryColor(ingredient.category)}`}>
            {ingredient.category}
          </span>
        </div>

        <div className="grid md:grid-cols-2 gap-6">
          <div>
            <h2 className="text-xl font-semibold mb-3">食材情報</h2>
            <div className="space-y-2">
              <div className="flex justify-between">
                <span className="text-gray-600">カテゴリ:</span>
                <span className="font-medium">{ingredient.category}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">カロリー:</span>
                <span className="font-medium">{ingredient.calories_per_100g} kcal/100g</span>
              </div>
            </div>
          </div>

          <div>
            <h2 className="text-xl font-semibold mb-3">説明</h2>
            <p className="text-gray-700">{ingredient.description}</p>
          </div>
        </div>
      </div>

      <div>
        <h2 className="text-2xl font-semibold mb-6">
          この食材を使用するレシピ ({recipes?.length}件)
        </h2>

        {!recipes || recipes?.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            この食材を使用するレシピが見つかりませんでした
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {recipes?.map(recipe => (
              <Link href={`/recipes/${recipe.id}`} key={recipe.id}>
                <div className="border rounded-lg p-4 shadow hover:shadow-lg transition-shadow cursor-pointer">
                  <h3 className="text-lg font-semibold mb-2">{recipe.name}</h3>
                  <p className="text-gray-600 text-sm mb-3 line-clamp-2">{recipe.description}</p>
                  
                  <div className="space-y-1 text-sm text-gray-600 mb-3">
                    <p>カテゴリ: {recipe.category}</p>
                    <p>準備: {recipe.prep_time_minutes}分 | 調理: {recipe.cook_time_minutes}分</p>
                    <p>{recipe.servings}人分</p>
                  </div>

                  <div className="flex justify-between items-center">
                    <span className={`px-2 py-1 rounded text-xs ${
                      recipe.difficulty === 'easy' ? 'bg-green-100 text-green-800' :
                      recipe.difficulty === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                      'bg-red-100 text-red-800'
                    }`}>
                      {recipe.difficulty === 'easy' ? '簡単' : 
                       recipe.difficulty === 'medium' ? '普通' : '難しい'}
                    </span>
                    <span className="text-xs text-gray-500">
                      合計{recipe.prep_time_minutes + recipe.cook_time_minutes}分
                    </span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>

      <div className="mt-8 text-center">
        <Link 
          href="/ingredients" 
          className="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          食材一覧に戻る
        </Link>
      </div>
    </div>
  );
}