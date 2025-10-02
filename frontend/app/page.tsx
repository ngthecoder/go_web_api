'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import LikeButton from '@/components/LikeButton';
import { useAuth } from '@/contexts/AuthContext';

interface Recipe {
  id: number;
  name: string;
  category: string;
  prep_time_minutes: number;
  cook_time_minutes: number;
  servings: number;
  difficulty: string;
  description: string;
  is_liked: boolean;
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
  const { token } = useAuth();
  
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
      
      const headers: HeadersInit = {};
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
      
      const response = await fetch(`http://localhost:8000/api/recipes?${params}`, {
        headers
      });
      const data = await response.json();
      setRecipeData(data);
    } catch (error) {
      console.error('API Error:', error);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchRecipes();
  }, [searchTerm, selectedCategory, selectedDifficulty, maxTime, sortBy, sortOrder, currentPage, itemsPerPage, token]);

  const handleLikeChange = (recipeId: number, isLiked: boolean) => {
    setRecipeData(prev => ({
      ...prev,
      recipes: prev.recipes.map(recipe => 
        recipe.id === recipeId 
          ? { ...recipe, is_liked: isLiked }
          : recipe
      )
    }));
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6">Recipes</h1>

      <div className="mb-6 space-y-4 bg-gray-50 p-6 rounded-lg">
        <div>
          <input
            type="text"
            placeholder="Search recipes..."
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
            <option value="">All Categories</option>
            <option value="Breakfast">Breakfast</option>
            <option value="Lunch">Lunch</option>
            <option value="Dinner">Dinner</option>
            <option value="Side">Side</option>
            <option value="Snack">Snack</option>
          </select>

          <select
            title="difficulty-selection"
            value={selectedDifficulty}
            onChange={(e) => setSelectedDifficulty(e.target.value)}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="">All Difficulties</option>
            <option value="easy">Easy</option>
            <option value="medium">Medium</option>
            <option value="hard">Hard</option>
          </select>

          <input
            type="number"
            placeholder="Max time (minutes)"
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
            <option value="6">Show 6</option>
            <option value="12">Show 12</option>
            <option value="24">Show 24</option>
            <option value="50">Show 50</option>
          </select>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <select
            title="sort-selection"
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="name">Sort by Name</option>
            <option value="prep_time">Sort by Prep Time</option>
            <option value="cook_time">Sort by Cook Time</option>
            <option value="total_time">Sort by Total Time</option>
            <option value="servings">Sort by Servings</option>
            <option value="difficulty">Sort by Difficulty</option>
          </select>

          <select
            title="order-selection"
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value)}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="asc">Ascending</option>
            <option value="desc">Descending</option>
          </select>
        </div>

        <div className="text-sm text-gray-600">
          Showing {((currentPage - 1) * itemsPerPage) + 1}-{Math.min(currentPage * itemsPerPage, recipeData.total)} of {recipeData.total} items
        </div>
      </div>

      {loading ? (
        <div className="text-center py-8">Loading...</div>
      ) : (
        <div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {recipeData.recipes?.map(recipe => (
              <div key={recipe.id} className="border rounded-lg p-4 shadow hover:shadow-lg transition-shadow relative">
                <div className="absolute top-4 right-4 z-10">
                  <LikeButton 
                    recipeId={recipe.id}
                    recipeName={recipe.name}
                    initialLiked={recipe.is_liked}
                    size="medium"
                    onLikeChange={(isLiked) => handleLikeChange(recipe.id, isLiked)}
                  />
                </div>

                <Link href={`/recipes/${recipe.id}`}>
                  <div className="cursor-pointer">
                    <h3 className="text-xl font-semibold mb-2 pr-10">{recipe.name}</h3>
                    <p className="text-gray-600 mb-2 text-sm line-clamp-2">{recipe.description}</p>
                    <div className="space-y-1 text-sm text-gray-600">
                      <p>Category: {recipe.category}</p>
                      <p>Prep: {recipe.prep_time_minutes}min | Cook: {recipe.cook_time_minutes}min</p>
                      <p>Serves {recipe.servings}</p>
                    </div>
                    <div className="mt-3 flex justify-between items-center">
                      <span className={`inline-block px-2 py-1 rounded text-sm ${
                        recipe.difficulty === 'easy' ? 'bg-green-100 text-green-800' :
                        recipe.difficulty === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                        'bg-red-100 text-red-800'
                      }`}>
                        {recipe.difficulty}
                      </span>
                      <span className="text-sm text-gray-500">
                        Total {recipe.prep_time_minutes + recipe.cook_time_minutes}min
                      </span>
                    </div>
                  </div>
                </Link>
              </div>
            ))}
          </div>

          {recipeData.total_pages > 1 && (
            <div className="mt-8 flex justify-center items-center space-x-2">
              <button
                onClick={() => setCurrentPage(currentPage - 1)}
                disabled={currentPage === 1}
                className="px-4 py-2 border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                Previous
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
                Next
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
