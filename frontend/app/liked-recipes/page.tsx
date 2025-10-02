'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import ProtectedRoute from '@/components/ProtectedRoute';
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
  instructions: string;
  description: string;
}

export default function LikedRecipesPage() {
  const { token } = useAuth();
  const [likedRecipes, setLikedRecipes] = useState<Recipe[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchLikedRecipes = async () => {
      if (!token) {
        setLoading(false);
        return;
      }

      try {
        const response = await fetch('http://localhost:8000/api/user/liked-recipes', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (!response.ok) {
          throw new Error('Failed to fetch liked recipes');
        }

        const data = await response.json();
        setLikedRecipes(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Something went wrong');
        console.error('Error fetching liked recipes:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchLikedRecipes();
  }, [token]);

  const handleLikeChange = (recipeId: number, isLiked: boolean) => {
    if (!isLiked) {
      setLikedRecipes(prev => prev.filter(recipe => recipe.id !== recipeId));
    }
  };

  return (
    <ProtectedRoute>
      <div className="container mx-auto p-4">
        <div className="mb-6">
          <h1 className="text-3xl font-bold mb-2">My Liked Recipes</h1>
          <p className="text-gray-600">
            {likedRecipes !== undefined && loading ? 'Loading...' : `You have ${likedRecipes?.length} liked recipe${likedRecipes?.length !== 1 ? 's' : ''}`}
          </p>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        {loading ? (
          <div className="flex justify-center items-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        ) : likedRecipes?.length === 0 ? (
          <div className="bg-white rounded-lg shadow-sm border p-12 text-center">
            <svg className="mx-auto h-16 w-16 text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
            </svg>
            <h3 className="text-xl font-medium text-gray-900 mb-2">No Liked Recipes Yet</h3>
            <p className="text-gray-600 mb-6">
              Start exploring recipes and click the heart icon to save your favorites!
            </p>
            <Link 
              href="/"
              className="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Browse Recipes
            </Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {likedRecipes?.map(recipe => (
              <div key={recipe.id} className="border rounded-lg p-4 shadow hover:shadow-lg transition-shadow relative">
                <div className="absolute top-4 right-4 z-10">
                  <LikeButton 
                    recipeId={recipe.id}
                    recipeName={recipe.name}
                    initialLiked={true}
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
        )}
      </div>
    </ProtectedRoute>
  );
}
