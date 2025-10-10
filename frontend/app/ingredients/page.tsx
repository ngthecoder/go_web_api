'use client';
import { useState, useEffect } from 'react';
import Link from 'next/link';
import { API_ENDPOINTS } from '@/lib/api-config';

interface Ingredient {
  id: number;
  name: string;
  category: string;
  calories_per_100g: number;
  description: string;
}

interface IngredientResponse {
  ingredients: Ingredient[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
  has_next: boolean;
}

export default function IngredientsPage() {
  const [ingredientData, setIngredientData] = useState<IngredientResponse>({
    ingredients: [],
    total: 0,
    page: 1,
    page_size: 20,
    total_pages: 0,
    has_next: false
  });
  const [loading, setLoading] = useState(true);
  
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');
  const [sortBy, setSortBy] = useState('name');
  const [sortOrder, setSortOrder] = useState('asc');
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(20);

  const fetchIngredients = async () => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      if (searchTerm) params.set('search', searchTerm);
      if (selectedCategory) params.set('category', selectedCategory);
      if (sortBy) params.set('sort', sortBy);
      if (sortOrder) params.set('order', sortOrder);
      params.set('page', currentPage.toString());
      params.set('limit', itemsPerPage.toString());
      
      const response = await fetch(`${API_ENDPOINTS.ingredients}?${params}`);
      const data = await response.json();
      setIngredientData(data);
    } catch (error) {
      console.error('API Error:', error);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchIngredients();
  }, [searchTerm, selectedCategory, sortBy, sortOrder, currentPage, itemsPerPage]);

  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
    window.scrollTo(0, 0);
  };

  const getCategoryColor = (category: string) => {
    const colors: Record<string, string> = {
      'Vegetables': 'bg-green-100 text-green-800',
      'Grains': 'bg-yellow-100 text-yellow-800',
      'Protein': 'bg-red-100 text-red-800',
      'Dairy': 'bg-blue-100 text-blue-800',
      'Seasonings': 'bg-purple-100 text-purple-800',
      'Fruits': 'bg-pink-100 text-pink-800',
      'Herbs & Spices': 'bg-indigo-100 text-indigo-800',
      'Pantry': 'bg-orange-100 text-orange-800',
    };
    return colors[category] || 'bg-gray-100 text-gray-800';
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6">Ingredients</h1>

      <div className="mb-6 space-y-4 bg-gray-50 p-6 rounded-lg">
        <div>
          <input
            type="text"
            placeholder="Search ingredients..."
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
            <option value="Vegetables">Vegetables</option>
            <option value="Grains">Grains</option>
            <option value="Protein">Protein</option>
            <option value="Dairy">Dairy</option>
            <option value="Seasonings">Seasonings</option>
            <option value="Fruits">Fruits</option>
            <option value="Herbs & Spices">Herbs & Spices</option>
            <option value="Pantry">Pantry</option>
          </select>

          <select
            title="sort-selection"
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="name">Sort by Name</option>
            <option value="calories">Sort by Calories</option>
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

          <select
            title="limit-selection"
            value={itemsPerPage}
            onChange={(e) => {
              setItemsPerPage(Number(e.target.value));
              setCurrentPage(1);
            }}
            className="px-4 py-2 border rounded-lg"
          >
            <option value="20">Show 20</option>
            <option value="40">Show 40</option>
            <option value="60">Show 60</option>
            <option value="100">Show 100</option>
          </select>
        </div>

        <div className="text-sm text-gray-600">
          Showing {((currentPage - 1) * itemsPerPage) + 1}-{Math.min(currentPage * itemsPerPage, ingredientData.total)} of {ingredientData.total} items
        </div>
      </div>

      {loading ? (
        <div className="text-center py-8">Loading...</div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {ingredientData?.ingredients?.map(ingredient => (
              <Link href={`/ingredients/${ingredient.id}`} key={ingredient.id}>
                <div className="border rounded-lg p-4 shadow hover:shadow-lg transition-shadow cursor-pointer">
                  <div className="flex justify-between items-start mb-2">
                    <h3 className="text-lg font-semibold">{ingredient.name}</h3>
                    <span className={`px-2 py-1 rounded text-xs ${getCategoryColor(ingredient.category)}`}>
                      {ingredient.category}
                    </span>
                  </div>
                  
                  <p className="text-gray-600 text-sm mb-3 line-clamp-2">{ingredient.description}</p>
                  
                  <div className="text-sm text-gray-500">
                    <span className="font-medium">{ingredient.calories_per_100g}</span> kcal/100g
                  </div>
                </div>
              </Link>
            ))}
          </div>

          {ingredientData?.total_pages > 1 && (
            <div className="mt-8 flex justify-center items-center space-x-2">
              <button
                onClick={() => handlePageChange(currentPage - 1)}
                disabled={currentPage === 1}
                className="px-4 py-2 border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                Previous
              </button>
              
              <div className="flex space-x-1">
                {Array.from({ length: Math.min(ingredientData?.total_pages, 5) }, (_, i) => {
                  const pageNum = Math.max(1, currentPage - 2) + i;
                  if (pageNum > ingredientData?.total_pages) return null;
                  
                  return (
                    <button
                      key={pageNum}
                      onClick={() => handlePageChange(pageNum)}
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
                onClick={() => handlePageChange(currentPage + 1)}
                disabled={!ingredientData.has_next}
                className="px-4 py-2 border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                Next
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}