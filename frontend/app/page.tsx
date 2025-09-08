'use client';
import { useState, useEffect } from 'react';
import Link from 'next/link';

interface Recipe {
  id: number;
  name: string;
  category: string;
  prep_time_minutes: number;
  cook_time_minutes: number;
  servings: number;
  difficulty: string;
  description: string;
}

interface RecipeResponse {
  recipes: Recipe[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
  has_next: boolean;
}

export default function HomePage() {
  const [recipeData, setRecipeData] = useState<RecipeResponse>({
    recipes: [],
    total: 0,
    page: 1,
    page_size: 10,
    total_pages: 0,
    has_next: false
  });
  const [loading, setLoading] = useState(true);

  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');
  const [selectedDifficulty, setSelectedDifficulty] = useState('');
  const [maxTime, setMaxTime] = useState('');
  const [sortBy, setSortBy] = useState('name');
  const [sortOrder, setSortOrder] = useState('asc');
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(12);

  const fetchRecipes = async () => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      if (searchTerm) params.set('search', searchTerm);
      if (selectedCategory) params.set('category', selectedCategory);
      if (selectedDifficulty) params.set('difficulty', selectedDifficulty);
      if (maxTime) params.set('max_time', maxTime);
      if (sortBy) params.set('sort', sortBy);
      if (sortOrder) params.set('order', sortOrder);
      params.set('page', currentPage.toString());
      params.set('limit', itemsPerPage.toString());
      
      const response = await fetch(`http://localhost:8000/api/recipes?${params}`);
      const data = await response.json();
      setRecipeData(data);
    } catch (error) {
      console.error('API Error:', error);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchRecipes();
  }, [searchTerm, selectedCategory, selectedDifficulty, maxTime, sortBy, sortOrder, currentPage, itemsPerPage]);

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6">レシピ一覧</h1>

      <div className="mb-6 space-y-4 bg-gray-50 p-6 rounded-lg">
        <div>
          <input
            type="text"
            placeholder="レシピを検索..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full px-4 py-2 border rounded-lg"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <select
            title="category-selection"
            value={selectedCategory}
            onChange={(e) => setSelectedCategory(e.target.value)}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="">すべてのカテゴリ</option>
            <option value="朝食">朝食</option>
            <option value="昼食">昼食</option>
            <option value="夕食">夕食</option>
            <option value="副菜">副菜</option>
            <option value="おやつ">おやつ</option>
          </select>

          <select
            title="difficulty-selection"
            value={selectedDifficulty}
            onChange={(e) => setSelectedDifficulty(e.target.value)}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="">すべての難易度</option>
            <option value="easy">簡単</option>
            <option value="medium">普通</option>
            <option value="hard">難しい</option>
          </select>

          <input
            type="number"
            placeholder="最大調理時間(分)"
            value={maxTime}
            onChange={(e) => setMaxTime(e.target.value)}
            className="px-4 py-2 border rounded-lg"
            min="1"
          />

          <select
            title="limit-selection"
            value={itemsPerPage}
            onChange={(e) => {
              setItemsPerPage(Number(e.target.value));
              setCurrentPage(1);
            }}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="6">6件表示</option>
            <option value="12">12件表示</option>
            <option value="24">24件表示</option>
            <option value="50">50件表示</option>
          </select>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <select
            title="sort-selection"
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="name">名前順</option>
            <option value="prep_time">準備時間順</option>
            <option value="cook_time">調理時間順</option>
            <option value="total_time">合計時間順</option>
            <option value="servings">人数順</option>
            <option value="difficulty">難易度順</option>
          </select>

          <select
            title="order-selection"
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value)}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="asc">昇順</option>
            <option value="desc">降順</option>
          </select>
        </div>

        <div className="text-sm text-gray-600">
          {recipeData.total}件中 {((currentPage - 1) * itemsPerPage) + 1}-{Math.min(currentPage * itemsPerPage, recipeData.total)}件を表示
        </div>
      </div>

      {loading ? (
        <div className="text-center py-8">Loading...</div>
      ) : (
        <div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {recipeData.recipes?.map(recipe => (
              <Link href={`/recipes/${recipe.id}`} key={recipe.id}>
                <div className="border rounded-lg p-4 shadow hover:shadow-lg transition-shadow cursor-pointer">
                  <h3 className="text-xl font-semibold mb-2">{recipe.name}</h3>
                  <p className="text-gray-600 mb-2 text-sm line-clamp-2">{recipe.description}</p>
                  <div className="space-y-1 text-sm text-gray-600">
                    <p>カテゴリ: {recipe.category}</p>
                    <p>準備: {recipe.prep_time_minutes}分 | 調理: {recipe.cook_time_minutes}分</p>
                    <p>{recipe.servings}人分</p>
                  </div>
                  <div className="mt-3 flex justify-between items-center">
                    <span className={`inline-block px-2 py-1 rounded text-sm ${
                      recipe.difficulty === 'easy' ? 'bg-green-100 text-green-800' :
                      recipe.difficulty === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                      'bg-red-100 text-red-800'
                    }`}>
                      {recipe.difficulty === 'easy' ? '簡単' : 
                       recipe.difficulty === 'medium' ? '普通' : '難しい'}
                    </span>
                    <span className="text-sm text-gray-500">
                      合計{recipe.prep_time_minutes + recipe.cook_time_minutes}分
                    </span>
                  </div>
                </div>
              </Link>
            ))}
          </div>

          {recipeData.total_pages > 1 && (
            <div className="mt-8 flex justify-center items-center space-x-2">
              <button
                onClick={() => setCurrentPage(currentPage - 1)}
                disabled={currentPage === 1}
                className="px-4 py-2 border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                前のページ
              </button>
              
              <div className="flex space-x-1">
                {Array.from({ length: Math.min(recipeData.total_pages, 5) }, (_, i) => {
                  const pageNum = Math.max(1, currentPage - 2) + i;
                  if (pageNum > recipeData.total_pages) return null;
                  
                  return (
                    <button
                      key={pageNum}
                      onClick={() => setCurrentPage(pageNum)}
                      className={`px-3 py-2 border rounded ${
                        currentPage === pageNum 
                          ? 'bg-blue-600 text-white' 
                          : 'hover:bg-gray-50'
                      }`}
                    >
                      {pageNum}
                    </button>
                  );
                })}
              </div>

              <button
                onClick={() => setCurrentPage(currentPage + 1)}
                disabled={!recipeData.has_next}
                className="px-4 py-2 border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                次のページ
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}