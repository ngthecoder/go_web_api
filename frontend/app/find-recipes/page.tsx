'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import LikeButton from '@/components/LikeButton';
import { useAuth } from '@/contexts/AuthContext';
import { API_ENDPOINTS } from '@/lib/api-config';

interface Ingredient {
  id: number;
  name: string;
  category: string;
}

interface MatchedRecipe {
  id: number;
  name: string;
  category: string;
  prep_time_minutes: number;
  cook_time_minutes: number;
  difficulty: string;
  matched_ingredients_count: number;
  total_ingredients_count: number;
  match_score: number;
  is_liked: boolean;
}

export default function FindRecipesPage() {
  const { token } = useAuth();
  const [ingredients, setIngredients] = useState<Ingredient[]>([]);
  const [selectedIngredients, setSelectedIngredients] = useState<number[]>([]);
  const [matchedRecipes, setMatchedRecipes] = useState<MatchedRecipe[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    fetch(`${API_ENDPOINTS.ingredients}?limit=100`)
      .then(res => res.json())
      .then(data => setIngredients(data.ingredients))
      .catch(err => console.error('Error fetching ingredients:', err));
  }, []);

  const toggleIngredient = (ingredientId: number) => {
    setSelectedIngredients(prev => 
      prev.includes(ingredientId)
        ? prev.filter(id => id !== ingredientId)
        : [...prev, ingredientId]
    );
  };

  const findRecipes = async () => {
    if (selectedIngredients.length === 0) return;
    
    setLoading(true);
    try {
      const ingredientIds = selectedIngredients.join(',');

      const headers: HeadersInit = {};
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
      
      const response = await fetch(
        `${API_ENDPOINTS.recipesByIngredients}?ingredients=${ingredientIds}&match_type=partial&limit=20`,
        { headers }
      );
      const data = await response.json();
      setMatchedRecipes(data);
    } catch (error) {
      console.error('Error finding recipes:', error);
    }
    setLoading(false);
  };

  const filteredIngredients = ingredients.filter(ingredient =>
    ingredient.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getSelectedIngredientNames = () => {
    return selectedIngredients
      .map(id => ingredients.find(ing => ing.id === id)?.name)
      .filter(Boolean);
  };

  const handleLikeChange = (recipeId: number, isLiked: boolean) => {
    setMatchedRecipes(prev => 
      prev.map(recipe => 
        recipe.id === recipeId 
          ? { ...recipe, is_liked: isLiked }
          : recipe
      )
    );
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6">Find Recipes by Ingredients</h1>
      
      <div className="grid lg:grid-cols-2 gap-8">
        <div>
          <h2 className="text-2xl font-semibold mb-4">Select Ingredients</h2>
          <input
            type="text"
            placeholder="Search ingredients..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full px-4 py-2 border rounded-lg mb-4"
          />

          {selectedIngredients.length > 0 && (
            <div className="mb-4 p-4 bg-blue-50 rounded-lg">
              <h3 className="font-semibold mb-2">Selected Ingredients ({selectedIngredients.length} items):</h3>
              <div className="flex flex-wrap gap-2">
                {getSelectedIngredientNames().map((name, index) => (
                  <span key={index} className="px-2 py-1 bg-blue-200 rounded text-sm">
                    {name}
                  </span>
                ))}
              </div>
              <button
                onClick={findRecipes}
                disabled={loading}
                className="mt-3 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? 'Searching Recipes...' : 'Find Recipes'}
              </button>
            </div>
          )}

          <div className="grid grid-cols-2 gap-2 max-h-96 overflow-y-auto">
            {filteredIngredients.map(ingredient => (
              <label
                key={ingredient.id}
                className={`flex items-center p-3 border rounded cursor-pointer transition-colors ${
                  selectedIngredients.includes(ingredient.id)
                    ? 'bg-blue-100 border-blue-500'
                    : 'hover:bg-gray-50'
                }`}
              >
                <input
                  type="checkbox"
                  checked={selectedIngredients.includes(ingredient.id)}
                  onChange={() => toggleIngredient(ingredient.id)}
                  className="mr-2"
                />
                <div>
                  <div className="font-medium">{ingredient.name}</div>
                  <div className="text-sm text-gray-500">{ingredient.category}</div>
                </div>
              </label>
            ))}
          </div>
        </div>

        <div>
          <h2 className="text-2xl font-semibold mb-4">
            Matched Recipes 
            {matchedRecipes.length > 0 && ` (${matchedRecipes.length} found)`}
          </h2>

          {matchedRecipes.length === 0 && !loading && selectedIngredients.length > 0 && (
            <div className="text-center py-8 text-gray-500">
              No recipes found with selected ingredients
            </div>
          )}

          {matchedRecipes.length === 0 && selectedIngredients.length === 0 && (
            <div className="text-center py-8 text-gray-500">
              Select ingredients to search for recipes
            </div>
          )}

          <div className="space-y-4">
            {matchedRecipes.map(recipe => (
              <div key={recipe.id} className="border rounded-lg p-4 hover:shadow-lg transition-shadow relative">
                <div className="absolute top-4 right-4">
                  <LikeButton 
                    recipeId={recipe.id}
                    recipeName={recipe.name}
                    initialLiked={recipe.is_liked}
                    size="medium"
                    onLikeChange={(isLiked) => handleLikeChange(recipe.id, isLiked)}
                  />
                </div>

                <Link href={`/recipes/${recipe.id}`}>
                  <div className="cursor-pointer pr-12">
                    <div className="flex justify-between items-start mb-2">
                      <h3 className="text-xl font-semibold">{recipe.name}</h3>
                    </div>
                    
                    <div className="mb-3">
                      <div className="flex items-center gap-2">
                        <div className="text-sm font-medium text-green-600">
                          Match: {Math.round(recipe.match_score * 100)}%
                        </div>
                        <div className="text-xs text-gray-500">
                          ({recipe.matched_ingredients_count}/{recipe.total_ingredients_count} ingredients)
                        </div>
                      </div>
                      
                      <div className="mt-2 w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-green-500 h-2 rounded-full transition-all"
                          style={{ width: `${recipe.match_score * 100}%` }}
                        ></div>
                      </div>
                    </div>
                    
                    <div className="flex justify-between items-center text-sm text-gray-600">
                      <span>{recipe.category}</span>
                      <span>{recipe.prep_time_minutes + recipe.cook_time_minutes} min</span>
                      <span className={`px-2 py-1 rounded ${
                        recipe.difficulty === 'easy' ? 'bg-green-100 text-green-800' :
                        recipe.difficulty === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                        'bg-red-100 text-red-800'
                      }`}>
                        {recipe.difficulty}
                      </span>
                    </div>
                  </div>
                </Link>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}