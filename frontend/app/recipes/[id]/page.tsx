'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import LikeButton from '@/components/LikeButton';
import { useAuth } from '@/contexts/AuthContext';
import { API_ENDPOINTS } from '@/lib/api-config';

interface IngredientWithQuantity {
  ingredient_id: number;
  name: string;
  quantity: number;
  unit: string;
  notes: string;
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
  is_liked: boolean;  // NEW
}

interface RecipeDetail {
  recipe: Recipe;
  ingredients: IngredientWithQuantity[];
}

export default function RecipeDetailPage() {
  const params = useParams();
  const { token } = useAuth();
  const [recipeDetail, setRecipeDetail] = useState<RecipeDetail | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (params.id) {
      const headers: HeadersInit = {};
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      fetch(API_ENDPOINTS.recipeById(String(params.id)), { headers })
        .then(res => res.json())
        .then(data => {
          setRecipeDetail(data);
          setLoading(false);
        })
        .catch(err => {
          console.error('API Error:', err);
          setLoading(false);
        });
    }
  }, [params.id, token]);

  const handleLikeChange = (isLiked: boolean) => {
    if (recipeDetail) {
      setRecipeDetail({
        ...recipeDetail,
        recipe: {
          ...recipeDetail.recipe,
          is_liked: isLiked
        }
      });
    }
  };

  if (loading) return <div className="container mx-auto p-4">Loading...</div>;
  if (!recipeDetail) return <div className="container mx-auto p-4">Recipe not found</div>;

  const { recipe, ingredients } = recipeDetail;

  return (
    <div className="container mx-auto p-4 max-w-4xl">
      <nav className="mb-6 text-sm text-gray-600">
        <Link href="/" className="hover:text-blue-600">Home</Link>
        <span className="mx-2">→</span>
        <span className="text-gray-900">{recipe.name}</span>
      </nav>

      <div className="bg-white rounded-lg shadow-sm border p-6 mb-6">
        <div className="flex justify-between items-start mb-4">
          <div className="flex-1">
            <h1 className="text-4xl font-bold mb-2">{recipe.name}</h1>
            <p className="text-gray-600">{recipe.description}</p>
          </div>
          
          <LikeButton 
            recipeId={recipe.id}
            recipeName={recipe.name}
            initialLiked={recipe.is_liked}
            size="large"
            showLabel
            onLikeChange={handleLikeChange}
          />
        </div>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-6">
          <div className="text-center p-3 bg-gray-50 rounded">
            <div className="text-2xl font-bold text-blue-600">{recipe.prep_time_minutes}</div>
            <div className="text-sm text-gray-600">Prep Time (min)</div>
          </div>
          <div className="text-center p-3 bg-gray-50 rounded">
            <div className="text-2xl font-bold text-green-600">{recipe.cook_time_minutes}</div>
            <div className="text-sm text-gray-600">Cook Time (min)</div>
          </div>
          <div className="text-center p-3 bg-gray-50 rounded">
            <div className="text-2xl font-bold text-purple-600">{recipe.servings}</div>
            <div className="text-sm text-gray-600">Servings</div>
          </div>
          <div className="text-center p-3 bg-gray-50 rounded">
            <div className={`inline-block px-3 py-1 rounded text-sm font-medium ${
              recipe.difficulty === 'easy' ? 'bg-green-100 text-green-800' :
              recipe.difficulty === 'medium' ? 'bg-yellow-100 text-yellow-800' :
              'bg-red-100 text-red-800'
            }`}>
              {recipe.difficulty}
            </div>
            <div className="text-sm text-gray-600 mt-1">Difficulty</div>
          </div>
        </div>
      </div>

      <div className="grid md:grid-cols-3 gap-6">
        <div className="md:col-span-1">
          <div className="bg-white rounded-lg shadow-sm border p-6 h-full flex flex-col">
            <h2 className="text-2xl font-semibold mb-4">Ingredients</h2>
            <ul className="space-y-3 flex-1">
              {ingredients?.map(ing => (
                <li key={ing.ingredient_id} className="flex justify-between items-start">
                  <Link 
                    href={`/ingredients/${ing.ingredient_id}`}
                    className="text-blue-600 hover:text-blue-800 hover:underline flex-1"
                  >
                    {ing.name}
                  </Link>
                  <span className="text-gray-600 text-sm ml-2 whitespace-nowrap">
                    {ing.quantity} {ing.unit}
                  </span>
                </li>
              ))}
            </ul>

            {ingredients?.length > 0 && (
              <Link 
                href={`/recipes/shopping-list/${recipe.id}`}
                className="mt-6 block w-full text-center px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
              >
                Create Shopping List
              </Link>
            )}
          </div>
        </div>
        
        <div className="md:col-span-2">
          <div className="bg-white rounded-lg shadow-sm border p-6 h-full">
            <h2 className="text-2xl font-semibold mb-4">Instructions</h2>
            <div className="prose max-w-none">
              <pre className="whitespace-pre-wrap font-sans text-gray-700 leading-relaxed">
                {recipe.instructions}
              </pre>
            </div>
          </div>
        </div>
      </div>
      
      <div className="mt-8 text-center">
        <Link 
          href="/" 
          className="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          ← Back to Recipes
        </Link>
      </div>
    </div>
  );
}