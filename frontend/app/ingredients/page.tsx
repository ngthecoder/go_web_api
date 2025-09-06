// app/ingredients/page.tsx
'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';

// API response types
type Ingredient = {
  id: number;
  name: string;
  category: string;
  calories_per_100g: number;
  description: string;
};

type IngredientsListResponse = {
  has_next: boolean;
  ingredients: Ingredient[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
};

// APIから食材のリストを取得する関数
async function getIngredients(page: number): Promise<IngredientsListResponse> {
  const res = await fetch(`http://localhost:8000/api/ingredients?page=${page}`);
  
  if (!res.ok) {
    throw new Error('APIから食材リストの取得に失敗しました。');
  }
  
  return res.json();
}

export default function IngredientsPage() {
  const [data, setData] = useState<IngredientsListResponse | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchIngredients = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await getIngredients(page);
        setData(response);
      } catch (e) {
        setError((e as Error).message);
      } finally {
        setLoading(false);
      }
    };

    fetchIngredients();
  }, [page]); // Re-run the effect whenever the 'page' state changes

  const handleNextPage = () => {
    if (data?.has_next) {
      setPage(prevPage => prevPage + 1);
    }
  };

  const handlePreviousPage = () => {
    if (data && data.page > 1) {
      setPage(prevPage => prevPage - 1);
    }
  };

  if (loading) {
    return <div className="text-center p-6 text-xl font-semibold text-gray-700">読み込み中...</div>;
  }

  if (error) {
    return (
      <div className="flex justify-center items-center h-screen bg-red-50 text-red-700">
        <div className="p-8 rounded-xl shadow-md border border-red-200">
          <h2 className="text-2xl font-bold">エラー</h2>
          <p className="mt-4 text-lg">{error}</p>
        </div>
      </div>
    );
  }

  if (!data || data.ingredients.length === 0) {
    return (
      <div className="flex justify-center items-center h-screen bg-gray-50 text-gray-700">
        <p className="text-xl">食材が見つかりませんでした。</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="container mx-auto max-w-5xl">
        <h1 className="text-5xl font-extrabold text-center text-gray-900 mb-12">食材リスト</h1>
        <ul className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
          {data.ingredients.map((ingredient) => (
            <li key={ingredient.id} className="bg-white rounded-xl shadow-lg transition-transform transform hover:scale-105 hover:shadow-2xl">
              <Link href={`/ingredients/${ingredient.id}`} className="block p-6">
                <div className="text-gray-400 text-sm font-medium mb-1">{ingredient.category}</div>
                <h2 className="text-2xl font-bold text-gray-800">{ingredient.name}</h2>
                <p className="mt-2 text-gray-600 leading-snug">{ingredient.description}</p>
                <div className="mt-4 text-sm text-gray-500">
                  <span className="font-semibold">カロリー:</span> {ingredient.calories_per_100g} kcal/100g
                </div>
              </Link>
            </li>
          ))}
        </ul>
        
        {/* Pagination controls */}
        <div className="mt-12 flex justify-center items-center space-x-6">
          <button
            onClick={handlePreviousPage}
            disabled={data.page <= 1}
            className="px-6 py-3 bg-blue-600 text-white font-semibold rounded-full shadow-md transition-colors duration-200 hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            &lt; 前へ
          </button>
          <span className="text-xl font-bold text-gray-700">
            {data.page} / {data.total_pages}
          </span>
          <button
            onClick={handleNextPage}
            disabled={!data.has_next}
            className="px-6 py-3 bg-blue-600 text-white font-semibold rounded-full shadow-md transition-colors duration-200 hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            次へ &gt;
          </button>
        </div>
      </div>
    </div>
  );
}